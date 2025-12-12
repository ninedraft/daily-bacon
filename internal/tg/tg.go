package tg

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

// HTTPDoer executes HTTP requests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client for Telegram Bot API.
type Client struct {
	doer   HTTPDoer
	apiURL string
	token  string
}

// New creates Client with provided HTTPDoer.
func New(doer HTTPDoer) *Client {
	return &Client{
		doer:   doer,
		apiURL: "https://api.telegram.org",
		token:  os.Getenv("TELEGRAM_TOKEN"),
	}
}

// Img represents image to send.
type Img struct {
	Name   string
	Reader io.Reader
}

// MediaUpload represents a single attachment inside a Telegram media group.
type MediaUpload struct {
	Type        string
	FileName    string
	Reader      io.Reader
	ContentType string
	Caption     string
}

// SendMessage sends text message. Images are ignored in this simple implementation.
func (c *Client) SendMessage(ctx context.Context, chatID, msg string, _ ...Img) error {
	if c.doer == nil {
		return errors.New("nil HTTPDoer")
	}
	if c.token == "" {
		return errors.New("empty token")
	}
	data := url.Values{}
	data.Set("chat_id", chatID)
	data.Set("text", msg)
	u := fmt.Sprintf("%s/bot%s/sendMessage", c.apiURL, c.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("new request url=%s: %w", u, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

// SendMediaGroup uploads multiple files as an album / media group.
func (c *Client) SendMediaGroup(ctx context.Context, chatID string, uploads []MediaUpload) error {
	if len(uploads) == 0 {
		return errors.New("media group requires at least one upload")
	}
	if c.doer == nil {
		return errors.New("nil HTTPDoer")
	}
	if c.token == "" {
		return errors.New("empty token")
	}

	type mediaItem struct {
		Type    string `json:"type"`
		Media   string `json:"media"`
		Caption string `json:"caption,omitempty"`
	}

	items := make([]mediaItem, 0, len(uploads))
	for i, upload := range uploads {
		if upload.Reader == nil {
			return fmt.Errorf("upload %d has nil reader", i)
		}
		mediaType := upload.Type
		if mediaType == "" {
			mediaType = "document"
		}
		items = append(items, mediaItem{
			Type:    mediaType,
			Media:   fmt.Sprintf("attach://file%d", i),
			Caption: upload.Caption,
		})
	}

	mediaJSON, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("marshal media payload: %w", err)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("chat_id", chatID); err != nil {
		return fmt.Errorf("write chat_id field: %w", err)
	}
	if err := writer.WriteField("media", string(mediaJSON)); err != nil {
		return fmt.Errorf("write media field: %w", err)
	}

	for i, upload := range uploads {
		fileName := upload.FileName
		if fileName == "" {
			fileName = fmt.Sprintf("file%d", i)
		}
		header := textproto.MIMEHeader{}
		header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, fileName))

		header.Set("Content-Type", "application/octet-stream")
		if upload.ContentType != "" {
			header.Set("Content-Type", upload.ContentType)
		}

		part, err := writer.CreatePart(header)
		if err != nil {
			return fmt.Errorf("create part: %w", err)
		}
		if _, err := io.Copy(part, upload.Reader); err != nil {
			return fmt.Errorf("copy upload %d: %w", i, err)
		}
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("close multipart writer: %w", err)
	}

	u := fmt.Sprintf("%s/bot%s/sendMediaGroup", c.apiURL, c.token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return fmt.Errorf("new request url=%s: %w", u, err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, payload)
	}
	return nil
}
