package meteo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

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
	Latitude     float64
	Longitude    float64
	Hourly       string
	Daily        string
	StartDate    string
	EndDate      string
	Timezone     string
	ForecastDays int
	PastDays     int
}

// AirQuality fetches air quality data.
func (c *Client) AirQuality(ctx context.Context, p Params) (models.AirQualityResponse, error) {
	var out models.AirQualityResponse
	u, err := url.Parse(c.url)
	if err != nil {
		return out, err
	}
	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%f", p.Latitude))
	q.Set("longitude", fmt.Sprintf("%f", p.Longitude))
	if p.Hourly != "" {
		q.Set("hourly", p.Hourly)
	}
	if p.Daily != "" {
		q.Set("daily", p.Daily)
	}
	if p.StartDate != "" {
		q.Set("start_date", p.StartDate)
	}
	if p.EndDate != "" {
		q.Set("end_date", p.EndDate)
	}
	if p.Timezone != "" {
		q.Set("timezone", p.Timezone)
	}
	if p.ForecastDays > 0 {
		q.Set("forecast_days", fmt.Sprintf("%d", p.ForecastDays))
	}
	if p.PastDays > 0 {
		q.Set("past_days", fmt.Sprintf("%d", p.PastDays))
	}
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return out, err
	}
	err = c.http.DoJSON(ctx, req, &out)
	return out, err
}
