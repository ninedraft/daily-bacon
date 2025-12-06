package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const meteoTimeLayout = "2006-01-02T15:04"

// WeatherResponse models Open-Meteo forecast response when using current weather.
type WeatherResponse struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	GenerationTimeMS     float64 `json:"generationtime_ms"`
	UTCOffsetSeconds     int     `json:"utc_offset_seconds"`
	Timezone             string  `json:"timezone"`
	TimezoneAbbreviation string  `json:"timezone_abbreviation"`
	Elevation            float64 `json:"elevation"`

	Current      *WeatherCurrent      `json:"current,omitempty"`
	CurrentUnits *WeatherCurrentUnits `json:"current_units,omitempty"`
}

// WeatherCurrentUnits holds units for requested current variables.
type WeatherCurrentUnits struct {
	Time               string `json:"time"`
	Interval           string `json:"interval"`
	Temperature2M      string `json:"temperature_2m,omitempty"`
	RelativeHumidity2M string `json:"relative_humidity_2m,omitempty"`
	WindSpeed10M       string `json:"wind_speed_10m,omitempty"`
}

// WeatherCurrent holds current weather metrics.
type WeatherCurrent struct {
	Time               time.Time `json:"time"`
	Interval           int       `json:"interval"`
	Temperature2M      float64   `json:"temperature_2m,omitempty"`
	RelativeHumidity2M float64   `json:"relative_humidity_2m,omitempty"`
	WindSpeed10M       float64   `json:"wind_speed_10m,omitempty"`
}

// UnmarshalJSON converts Open-Meteo time strings into time.Time values.
func (w *WeatherCurrent) UnmarshalJSON(data []byte) error {
	type rawWeatherCurrent struct {
		Time               string  `json:"time"`
		Interval           int     `json:"interval"`
		Temperature2M      float64 `json:"temperature_2m,omitempty"`
		RelativeHumidity2M float64 `json:"relative_humidity_2m,omitempty"`
		WindSpeed10M       float64 `json:"wind_speed_10m,omitempty"`
	}

	var raw rawWeatherCurrent
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("decode weather current: %w", err)
	}

	var parsed time.Time
	if raw.Time != "" {
		var err error
		parsed, err = time.Parse(meteoTimeLayout, raw.Time)
		if err != nil {
			return fmt.Errorf("parse current time %q: %w", raw.Time, err)
		}
	}

	*w = WeatherCurrent{
		Time:               parsed,
		Interval:           raw.Interval,
		Temperature2M:      raw.Temperature2M,
		RelativeHumidity2M: raw.RelativeHumidity2M,
		WindSpeed10M:       raw.WindSpeed10M,
	}
	return nil
}
