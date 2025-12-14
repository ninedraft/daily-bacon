package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"mime/multipart"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/time/rate"

	"github.com/ninedraft/daily-bacon/internal/tg"

	_ "golang.org/x/crypto/x509roots/fallback"
)

const (
	envChatID      = "DAILY_BACON_CHAT_ID"
	envTokenFile   = "TELEGRAM_TOKEN_FILE"
	envGatewayAddr = "DAILY_BACON_GATEWAY_ADDR"
	envChatsConfig = "DAILY_BACON_CHAT_CONFIG"

	defaultGatewayAddr = ":8080"
	maxMultipartMemory = 64 << 20 // 64MB
	maxCaptionRunes    = 1024
)

type chatInfo struct {
	Label string
	ID    string
}

type chatResolverFunc func(*http.Request) (chatInfo, error)

type chatLookupError struct {
	Label  string
	Labels []string
}

func (e chatLookupError) Error() string {
	return fmt.Sprintf("chat label %q not found", e.Label)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("gateway failed to start", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	chatID := os.Getenv(envChatID)
	if chatID == "" {
		return fmt.Errorf("%s must be set", envChatID)
	}

	tokenFile := os.Getenv(envTokenFile)
	if tokenFile == "" {
		return fmt.Errorf("%s must be set", envTokenFile)
	}
	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		return fmt.Errorf("read token file: %w", err)
	}
	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return fmt.Errorf("token file %s is empty", tokenFile)
	}
	if err := os.Setenv("TELEGRAM_TOKEN", token); err != nil {
		return fmt.Errorf("set TELEGRAM_TOKEN: %w", err)
	}

	configPath := os.Getenv(envChatsConfig)
	if configPath == "" {
		return fmt.Errorf("missing configmap from %q", envChatsConfig)
	}

	chatMap, err := loadChatConfig(configPath)
	if err != nil {
		return fmt.Errorf("load chat config %s: %w", configPath, err)
	}

	client := tg.New(http.DefaultClient)

	addr := os.Getenv(envGatewayAddr)
	if addr == "" {
		addr = defaultGatewayAddr
	}

	limiter := rate.NewLimiter(rate.Every(time.Second/2), 1)

	mux := http.NewServeMux()
	defaultResolver := func(*http.Request) (chatInfo, error) {
		return chatInfo{Label: "default", ID: chatID}, nil
	}

	mux.HandleFunc("/message", messageHandler(logger, client, limiter, defaultResolver))
	mux.HandleFunc("/message/{label}", messageHandler(logger, client, limiter, newChatResolver(chatMap)))

	server := &http.Server{
		Addr:         addr,
		Handler:      mwLog(logger, mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Info("gateway starting", "addr", addr, "chat", chatID)
	return server.ListenAndServe()
}

func messageHandler(logger *slog.Logger, client *tg.Client, limiter *rate.Limiter, resolver chatResolverFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Wait(r.Context()); err != nil {
			slog.Error("waiting for limit", "error", err)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(maxMultipartMemory); err != nil {
			logger.Error("parse multipart form", "err", err)
			http.Error(w, "invalid multipart payload", http.StatusBadRequest)
			return
		}
		defer func() {
			if r.MultipartForm != nil {
				_ = r.MultipartForm.RemoveAll()
			}
		}()

		chat, err := resolver(r)
		if err != nil {
			if lookupErr, ok := err.(chatLookupError); ok {
				writeChatLookupError(w, lookupErr)
				return
			}
			logger.Error("resolve chat", "err", err)
			http.Error(w, "unable to resolve chat", http.StatusInternalServerError)
			return
		}

		text := strings.TrimSpace(r.FormValue("text"))
		captionText := text
		needsSeparateText := false
		if text != "" && utf8.RuneCountInString(text) > maxCaptionRunes {
			captionText = ""
			needsSeparateText = true
		}
		uploads, err := collectUploads(r.MultipartForm)
		if err != nil {
			logger.Error("collect uploads", "err", err)
			http.Error(w, "failed to read uploads", http.StatusBadRequest)
			return
		}

		logger.Info("incoming request", "remote", r.RemoteAddr, "files", len(uploads), "has_text", text != "", "chat_label", chat.Label)

		switch {
		case len(uploads) == 0 && text == "":
			logger.Info("nothing to send, skipping", "chat_label", chat.Label)
			w.WriteHeader(http.StatusNoContent)
			return
		case len(uploads) == 0:
			if err := client.SendMessage(r.Context(), chat.ID, text); err != nil {
				logger.Error("send text message", "err", err, "chat_label", chat.Label)
				http.Error(w, "failed to deliver message", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusAccepted)
			return
		}

		type fileEntry struct {
			ContentType string
		}
		logEntry := map[string]fileEntry{}

		media := make([]tg.MediaUpload, len(uploads))
		for i, upload := range uploads {
			var caption string
			if i == 0 {
				caption = captionText
			}
			media[i] = tg.MediaUpload{
				FileName:    upload.FileName,
				Reader:      upload.Reader,
				ContentType: upload.ContentType,
				Caption:     caption,
			}

			logEntry[upload.FileName] = fileEntry{
				ContentType: upload.ContentType,
			}
		}

		logger.Info("sending files", "files", logEntry, "chat_label", chat.Label)

		if err := client.SendMediaGroup(r.Context(), chat.ID, media); err != nil {
			logger.Error("send media group", "err", err, "chat_label", chat.Label)
			http.Error(w, "failed to deliver media group", http.StatusInternalServerError)
			return
		}

		if needsSeparateText {
			if err := client.SendMessage(r.Context(), chat.ID, text); err != nil {
				logger.Error("send text message after media group", "err", err, "chat_label", chat.Label)
				http.Error(w, "failed to deliver media group text", http.StatusInternalServerError)
				return
			}
		}

		logger.Info("delivered", "chat_label", chat.Label, "chat_id", chat.ID)
		w.WriteHeader(http.StatusAccepted)
	}
}

type collectedUpload struct {
	Reader      io.Reader
	FileName    string
	ContentType string
}

func collectUploads(form *multipart.Form) ([]collectedUpload, error) {
	if form == nil {
		return nil, nil
	}

	var uploads []collectedUpload

	handleFile := func(fileheader *multipart.FileHeader) (collectedUpload, error) {
		file, err := fileheader.Open()
		if err != nil {
			return collectedUpload{}, fmt.Errorf("open %s: %w", fileheader.Filename, err)
		}
		defer file.Close()

		re, detectedContentType, err := detectContentType(file)
		if err != nil {
			return collectedUpload{}, err
		}

		return collectedUpload{
			Reader:      re,
			FileName:    fileheader.Filename,
			ContentType: cmp.Or(fileheader.Header.Get("Content-Type"), detectedContentType),
		}, nil
	}

	for _, list := range form.File {
		for _, fh := range list {
			upload, err := handleFile(fh)
			if err != nil {
				return nil, fmt.Errorf("handling file %q: %w", fh.Filename, err)
			}
			uploads = append(uploads, upload)
		}
	}

	return uploads, nil
}

func detectContentType(re io.Reader) (_ io.Reader, contentType string, _ error) {
	head := make([]byte, 512)

	n, err := re.Read(head)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, "", err
	}
	head = head[:n]

	return io.MultiReader(bytes.NewReader(head), re), http.DetectContentType(head), nil
}

func mwLog(log *slog.Logger, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("request",
			"method", r.Method,
			"URL", r.URL)

		next.ServeHTTP(w, r)
	}
}

func newChatResolver(chats map[string]string) chatResolverFunc {
	labels := slices.Sorted(maps.Keys(chats))

	return func(r *http.Request) (chatInfo, error) {
		label := r.PathValue("label")
		if label == "" {
			return chatInfo{}, chatLookupError{Label: label, Labels: labels}
		}
		if chatID, ok := chats[label]; ok {
			return chatInfo{Label: label, ID: chatID}, nil
		}
		return chatInfo{}, chatLookupError{Label: label, Labels: labels}
	}
}

func writeChatLookupError(w http.ResponseWriter, err chatLookupError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":  err.Error(),
		"labels": err.Labels,
	})
}

func loadChatConfig(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg struct {
		Chats map[string]string `json:"chats"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Chats == nil {
		cfg.Chats = map[string]string{}
	}
	return cfg.Chats, nil
}
