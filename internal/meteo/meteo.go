package meteo

import (
	"context"
	"fmt"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/models"
)

const baseURL = "https://air-quality-api.open-meteo.com/v1/air-quality"

// Client for Open-Meteo air quality API.
type Client struct {
	http *client.Client
	url  string
}

// New creates Open-Meteo client.
func New(cl *client.Client) *Client {
	if cl == nil {
		cl = client.New(nil)
	}
	return &Client{http: cl, url: baseURL}
}

// Params for air quality request.
type Params = RequestParams

// AirQuality fetches air quality data.
func (c *Client) AirQuality(ctx context.Context, p Params) (models.AirQualityResponse, error) {
	var out models.AirQualityResponse

	req, err := newRequest(ctx, c.url, p)
	if err != nil {
		return out, err
	}

	err = c.http.DoJSON(ctx, req, &out)
	if err != nil {
		return out, fmt.Errorf("do request: %w", err)
	}
	return out, nil
}
