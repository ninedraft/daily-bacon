package models

import "time"

// AirQualityResponse models the complete JSON response.
type AirQualityResponse struct {
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	GenerationTimeMS float64 `json:"generationtime_ms"`
	UTCOffsetSeconds int     `json:"utc_offset_seconds"`
	Timezone         string  `json:"timezone"`

	Current      *CurrentData  `json:"current"`
	CurrentUnits *CurrentUnits `json:"current_units"`

	DailyUnits *DailyUnits `json:"daily_units,omitempty"`
	Daily      *DailyData  `json:"daily,omitempty"`
}

// HourlyUnits holds the unit labels for each hourly field.
type HourlyUnits struct {
	Time            string `json:"time"`
	PM10            string `json:"pm10,omitempty"`
	PM25            string `json:"pm2_5,omitempty"`
	Ozone           string `json:"ozone,omitempty"`
	NitrogenDioxide string `json:"nitrogen_dioxide,omitempty"`
	CarbonMonoxide  string `json:"carbon_monoxide,omitempty"`
	SulphurDioxide  string `json:"sulphur_dioxide,omitempty"`
}

// HourlyData holds the time series values.
type HourlyData struct {
	Time            []time.Time `json:"time"`
	PM10            []float64   `json:"pm10,omitempty"`
	PM25            []float64   `json:"pm2_5,omitempty"`
	Ozone           []float64   `json:"ozone,omitempty"`
	NitrogenDioxide []float64   `json:"nitrogen_dioxide,omitempty"`
	CarbonMonoxide  []float64   `json:"carbon_monoxide,omitempty"`
	SulphurDioxide  []float64   `json:"sulphur_dioxide,omitempty"`
}

type CurrentData struct {
	PM10            float64 `json:"pm10,omitempty"`
	PM25            float64 `json:"pm2_5,omitempty"`
	Ozone           float64 `json:"ozone,omitempty"`
	NitrogenDioxide float64 `json:"nitrogen_dioxide,omitempty"`
	CarbonMonoxide  float64 `json:"carbon_monoxide,omitempty"`
	SulphurDioxide  float64 `json:"sulphur_dioxide,omitempty"`
}

type CurrentUnits struct {
	PM10            string `json:"pm10,omitempty"`
	PM25            string `json:"pm2_5,omitempty"`
	Ozone           string `json:"ozone,omitempty"`
	NitrogenDioxide string `json:"nitrogen_dioxide,omitempty"`
	CarbonMonoxide  string `json:"carbon_monoxide,omitempty"`
	SulphurDioxide  string `json:"sulphur_dioxide,omitempty"`
}

// DailyUnits holds unit labels for each daily statistic.
type DailyUnits struct {
	Time     string `json:"time"`
	PM10Max  string `json:"pm10_max,omitempty"`
	PM10Mean string `json:"pm10_mean,omitempty"`
	PM10Min  string `json:"pm10_min,omitempty"`
}

// DailyData holds daily time series values.
type DailyData struct {
	PM10Max  []float64 `json:"pm10_max,omitempty"`
	PM10Mean []float64 `json:"pm10_mean,omitempty"`
	PM10Min  []float64 `json:"pm10_min,omitempty"`
}
