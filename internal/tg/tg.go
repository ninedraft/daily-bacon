package tg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// HTTPDoer executes HTTP requests.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
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
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.doer.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
