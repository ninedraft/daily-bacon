package meteo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

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
type Params struct {
	Latitude           float64
	Longitude          float64
	Current            []string
	Hourly             []string
	Daily              []string
	StartDate, EndDate time.Time
	Timezone           string
	ForecastDays       int
	PastDays           int
}

// AirQuality fetches air quality data.
func (c *Client) AirQuality(ctx context.Context, p Params) (models.AirQualityResponse, error) {
	var out models.AirQualityResponse
	u, err := url.Parse(c.url)
	if err != nil {
		return out, fmt.Errorf("parse url=%s: %w", c.url, err)
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%f", p.Latitude))
	q.Set("longitude", fmt.Sprintf("%f", p.Longitude))

	if len(p.Current) > 0 {
		q.Set("current", strings.Join(p.Current, ","))
	}
	if len(p.Hourly) > 0 {
		q.Set("hourly", strings.Join(p.Hourly, ","))
	}
	if len(p.Daily) > 0 {
		q.Set("daily", strings.Join(p.Daily, ","))
	}

	if !p.StartDate.IsZero() {
		q.Set("start_date", p.StartDate.Format(time.DateOnly))
	}
	if !p.EndDate.IsZero() {
		q.Set("end_date", p.EndDate.Format(time.DateOnly))
	}
	if p.Timezone != "" {
		q.Set("timezone", p.Timezone)
	}
	if p.ForecastDays > 0 {
		q.Set("forecast_days", strconv.Itoa(p.ForecastDays))
	}
	if p.PastDays > 0 {
		q.Set("past_days", strconv.Itoa(p.PastDays))
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return out, fmt.Errorf("new request url=%s: %w", u, err)
	}

	err = c.http.DoJSON(ctx, req, &out)
	if err != nil {
		return out, fmt.Errorf("do request: %w", err)
	}
	return out, nil
}
