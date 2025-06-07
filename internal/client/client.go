package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Client wraps http.Client.
type Client struct {
	httpClient *http.Client
}

// ErrBadResponseScheme is returned when request URL scheme is not HTTP or HTTPS.
var ErrBadResponseScheme = errors.New("bad response scheme")

// New creates a Client. Transport may be nil.
func New(transport http.RoundTripper) *Client {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &Client{httpClient: &http.Client{Transport: transport}}
}

// ErrUnexpectedStatus is returned for non-2xx status codes.
type ErrUnexpectedStatus struct {
	Code int
	Body []byte
}

func (e *ErrUnexpectedStatus) Error() string {
	return fmt.Sprintf("unexpected status %d", e.Code)
}

// DoJSON sends the request and decodes JSON response into dst.
func (c *Client) DoJSON(ctx context.Context, req *http.Request, dst any) error {
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		return ErrBadResponseScheme
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &ErrUnexpectedStatus{Code: resp.StatusCode, Body: body}
	}
	if dst != nil {
		if err := json.Unmarshal(body, dst); err != nil {
			return err
		}
	}
	return nil
}
