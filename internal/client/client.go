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

// UnexpectedStatusError is returned for non-2xx status codes.
type UnexpectedStatusError struct {
	Code int
	Body []byte
}

func (e *UnexpectedStatusError) Error() string {
	return fmt.Sprintf("unexpected status %d: %s", e.Code, e.Body)
}

const bodySizeLimit = 10_000_000

// DoJSON sends the request and decodes JSON response into dst.
func (c *Client) DoJSON(ctx context.Context, req *http.Request, dst any) error {
	if req.URL.Scheme != "http" && req.URL.Scheme != "https" {
		return ErrBadResponseScheme
	}
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request url=%s: %w", req.URL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, bodySizeLimit))
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &UnexpectedStatusError{Code: resp.StatusCode, Body: body}
	}
	if dst != nil {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("unmarshal json: %w", err)
		}
	}
	return nil
}
