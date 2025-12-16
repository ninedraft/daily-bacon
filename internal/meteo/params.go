package meteo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// RequestParams describe the shared request options accepted by Open-Meteo
// endpoints.
type RequestParams struct {
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

func (p RequestParams) apply(u *url.URL) {
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
}

func newRequest(ctx context.Context, baseURL string, p RequestParams) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse url=%s: %w", baseURL, err)
	}

	p.apply(u)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request url=%s: %w", u, err)
	}

	return req, nil
}
