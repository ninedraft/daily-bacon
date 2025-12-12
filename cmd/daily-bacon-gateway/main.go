package main

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
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

	defaultGatewayAddr = ":8080"
	maxMultipartMemory = 64 << 20 // 64MB
	maxCaptionRunes    = 1024
)

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

	client := tg.New(http.DefaultClient)

	addr := os.Getenv(envGatewayAddr)
	if addr == "" {
		addr = defaultGatewayAddr
	}

	limiter := rate.NewLimiter(rate.Every(time.Second/2), 1)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /message", messageHandler(logger, client, limiter, chatID))

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

func messageHandler(logger *slog.Logger, client *tg.Client, limiter *rate.Limiter, chatID string) http.HandlerFunc {
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

		text := strings.TrimSpace(r.FormValue("text"))
		captionText := text
		needsSeparateText := false
		if utf8.RuneCountInString(text) > maxCaptionRunes {
			captionText = ""
			needsSeparateText = true
		}
		uploads, err := collectUploads(r.MultipartForm)
		if err != nil {
			logger.Error("collect uploads", "err", err)
			http.Error(w, "failed to read uploads", http.StatusBadRequest)
			return
		}

		logger.Info("incoming request", "remote", r.RemoteAddr, "files", len(uploads), "has_text", text != "")

		switch {
		case len(uploads) == 0 && text == "":
			logger.Info("nothing to send, skipping")
			w.WriteHeader(http.StatusNoContent)
			return
		case len(uploads) == 0:
			if err := client.SendMessage(r.Context(), chatID, text); err != nil {
				logger.Error("send text message", "err", err)
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

		logger.Info("sending files", "files", logEntry)

		if err := client.SendMediaGroup(r.Context(), chatID, media); err != nil {
			logger.Error("send media group", "err", err)
			http.Error(w, "failed to deliver media group", http.StatusInternalServerError)
			return
		}

		if needsSeparateText {
			if err := client.SendMessage(r.Context(), chatID, text); err != nil {
				logger.Error("send text message after media group", "err", err)
				http.Error(w, "failed to deliver media group text", http.StatusInternalServerError)
				return
			}
		}

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
