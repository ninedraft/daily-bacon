package meteo

import (
	"context"
	"fmt"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/models"
)

const forecastBaseURL = "https://api.open-meteo.com/v1/forecast"

// WeatherClient requests forecast data from Open-Meteo.
type WeatherClient struct {
	http *client.Client
	url  string
}

// NewWeather creates Open-Meteo weather client.
func NewWeather(cl *client.Client) *WeatherClient {
	if cl == nil {
		cl = client.New(nil)
	}
	return &WeatherClient{http: cl, url: forecastBaseURL}
}

// WeatherParams describes how to build an Open-Meteo forecast request.
//
// The latitude and longitude are required and must be passed in decimal degrees
// (positive for northern/eastern hemispheres, negative for southern/western).
//
// The Current, Hourly, and Daily slices define which variables to request in
// each resolution. Use the constants in weather_vars.go (for example,
// Temperature2M or RelativeHumidity2M) or strings from the Open-Meteo API
// reference. Empty slices omit that section from the response entirely.
//
// StartDate and EndDate constrain the returned range to specific calendar days
// using the server timezone; leave them zero-valued to let the API determine
// the period based on ForecastDays and PastDays.
//
// Timezone sets the IANA timezone name used for time values in the response.
// If unset, Open-Meteo uses the best match for the provided coordinates.
//
// ForecastDays limits how many upcoming days are returned (up to the provider
// limit). PastDays requests historical data counting backward from today. Both
// are optional and ignored when zero.
type WeatherParams = RequestParams

// Forecast fetches weather data.
func (c *WeatherClient) Forecast(ctx context.Context, p WeatherParams) (models.WeatherResponse, error) {
	var out models.WeatherResponse

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
