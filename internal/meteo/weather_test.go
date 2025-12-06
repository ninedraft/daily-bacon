package meteo

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ninedraft/daily-bacon/internal/client"
)

func TestWeatherForecastNicosia(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	cl := client.New(http.DefaultClient.Transport)
	weather := NewWeather(cl)

	resp, err := weather.Forecast(ctx, WeatherParams{
		Latitude:  35.1856,
		Longitude: 33.3823,
		Current: []string{
			Temperature2M,
			RelativeHumidity2M,
			WindSpeed10M,
		},
		Timezone: "auto",
	})
	if err != nil {
		t.Fatalf("forecast request failed: %v", err)
	}

	if resp.Current == nil {
		t.Fatalf("missing current weather data")
	}
	if resp.CurrentUnits == nil {
		t.Fatalf("missing current units")
	}
	if resp.Current.Time.IsZero() {
		t.Fatalf("current timestamp is empty")
	}
	if resp.Timezone == "" {
		t.Fatalf("timezone is empty")
	}
}
