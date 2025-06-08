package models

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
	Time                string `json:"time"`
	PM10                string `json:"pm10,omitempty"`
	PM25                string `json:"pm2_5,omitempty"`
	CarbonMonoxide      string `json:"carbon_monoxide,omitempty"`
	CarbonDioxide       string `json:"carbon_dioxide,omitempty"`
	NitrogenDioxide     string `json:"nitrogen_dioxide,omitempty"`
	SulphurDioxide      string `json:"sulphur_dioxide,omitempty"`
	Ozone               string `json:"ozone,omitempty"`
	AerosolOpticalDepth string `json:"aerosol_optical_depth,omitempty"`
	Dust                string `json:"dust,omitempty"`
	UVIndex             string `json:"uv_index,omitempty"`
	UVIndexClearSky     string `json:"uv_index_clear_sky,omitempty"`
	Ammonia             string `json:"ammonia,omitempty"`
	Methane             string `json:"methane,omitempty"`
	AlderPollen         string `json:"alder_pollen,omitempty"`
	BirchPollen         string `json:"birch_pollen,omitempty"`
	GrassPollen         string `json:"grass_pollen,omitempty"`
	MugwortPollen       string `json:"mugwort_pollen,omitempty"`
	OlivePollen         string `json:"olive_pollen,omitempty"`
	RagweedPollen       string `json:"ragweed_pollen,omitempty"`
	EuropeanAQI         string `json:"european_aqi,omitempty"`
	EuropeanAQIPM25     string `json:"european_aqi_pm2_5,omitempty"`
	EuropeanAQIPM10     string `json:"european_aqi_pm10,omitempty"`
	EuropeanAQINO2      string `json:"european_aqi_nitrogen_dioxide,omitempty"`
	EuropeanAQIOzone    string `json:"european_aqi_ozone,omitempty"`
	EuropeanAQISO2      string `json:"european_aqi_sulphur_dioxide,omitempty"`
	USAQI               string `json:"us_aqi,omitempty"`
	USAQIPM25           string `json:"us_aqi_pm2_5,omitempty"`
	USAQIPM10           string `json:"us_aqi_pm10,omitempty"`
	USAQINO2            string `json:"us_aqi_nitrogen_dioxide,omitempty"`
	USAQIOzone          string `json:"us_aqi_ozone,omitempty"`
	USAQISO2            string `json:"us_aqi_sulphur_dioxide,omitempty"`
	USAQICarbonMonoxide string `json:"us_aqi_carbon_monoxide,omitempty"`
}

// HourlyData holds the time series values.
type HourlyData struct {
	PM10                []float64 `json:"pm10,omitempty"`
	PM25                []float64 `json:"pm2_5,omitempty"`
	CarbonMonoxide      []float64 `json:"carbon_monoxide,omitempty"`
	CarbonDioxide       []float64 `json:"carbon_dioxide,omitempty"`
	NitrogenDioxide     []float64 `json:"nitrogen_dioxide,omitempty"`
	SulphurDioxide      []float64 `json:"sulphur_dioxide,omitempty"`
	Ozone               []float64 `json:"ozone,omitempty"`
	AerosolOpticalDepth []float64 `json:"aerosol_optical_depth,omitempty"`
	Dust                []float64 `json:"dust,omitempty"`
	UVIndex             []float64 `json:"uv_index,omitempty"`
	UVIndexClearSky     []float64 `json:"uv_index_clear_sky,omitempty"`
	Ammonia             []float64 `json:"ammonia,omitempty"`
	Methane             []float64 `json:"methane,omitempty"`
	AlderPollen         []float64 `json:"alder_pollen,omitempty"`
	BirchPollen         []float64 `json:"birch_pollen,omitempty"`
	GrassPollen         []float64 `json:"grass_pollen,omitempty"`
	MugwortPollen       []float64 `json:"mugwort_pollen,omitempty"`
	OlivePollen         []float64 `json:"olive_pollen,omitempty"`
	RagweedPollen       []float64 `json:"ragweed_pollen,omitempty"`
	EuropeanAQI         []float64 `json:"european_aqi,omitempty"`
	EuropeanAQIPM25     []float64 `json:"european_aqi_pm2_5,omitempty"`
	EuropeanAQIPM10     []float64 `json:"european_aqi_pm10,omitempty"`
	EuropeanAQINO2      []float64 `json:"european_aqi_nitrogen_dioxide,omitempty"`
	EuropeanAQIOzone    []float64 `json:"european_aqi_ozone,omitempty"`
	EuropeanAQISO2      []float64 `json:"european_aqi_sulphur_dioxide,omitempty"`
	USAQI               []float64 `json:"us_aqi,omitempty"`
	USAQIPM25           []float64 `json:"us_aqi_pm2_5,omitempty"`
	USAQIPM10           []float64 `json:"us_aqi_pm10,omitempty"`
	USAQINO2            []float64 `json:"us_aqi_nitrogen_dioxide,omitempty"`
	USAQIOzone          []float64 `json:"us_aqi_ozone,omitempty"`
	USAQISO2            []float64 `json:"us_aqi_sulphur_dioxide,omitempty"`
	USAQICarbonMonoxide []float64 `json:"us_aqi_carbon_monoxide,omitempty"`
}

type CurrentData struct {
	PM10                float64 `json:"pm10,omitempty"`
	PM25                float64 `json:"pm2_5,omitempty"`
	CarbonMonoxide      float64 `json:"carbon_monoxide,omitempty"`
	CarbonDioxide       float64 `json:"carbon_dioxide,omitempty"`
	NitrogenDioxide     float64 `json:"nitrogen_dioxide,omitempty"`
	SulphurDioxide      float64 `json:"sulphur_dioxide,omitempty"`
	Ozone               float64 `json:"ozone,omitempty"`
	AerosolOpticalDepth float64 `json:"aerosol_optical_depth,omitempty"`
	Dust                float64 `json:"dust,omitempty"`
	UVIndex             float64 `json:"uv_index,omitempty"`
	UVIndexClearSky     float64 `json:"uv_index_clear_sky,omitempty"`
	Ammonia             float64 `json:"ammonia,omitempty"`
	Methane             float64 `json:"methane,omitempty"`
	AlderPollen         float64 `json:"alder_pollen,omitempty"`
	BirchPollen         float64 `json:"birch_pollen,omitempty"`
	GrassPollen         float64 `json:"grass_pollen,omitempty"`
	MugwortPollen       float64 `json:"mugwort_pollen,omitempty"`
	OlivePollen         float64 `json:"olive_pollen,omitempty"`
	RagweedPollen       float64 `json:"ragweed_pollen,omitempty"`
	EuropeanAQI         float64 `json:"european_aqi,omitempty"`
	EuropeanAQIPM25     float64 `json:"european_aqi_pm2_5,omitempty"`
	EuropeanAQIPM10     float64 `json:"european_aqi_pm10,omitempty"`
	EuropeanAQINO2      float64 `json:"european_aqi_nitrogen_dioxide,omitempty"`
	EuropeanAQIOzone    float64 `json:"european_aqi_ozone,omitempty"`
	EuropeanAQISO2      float64 `json:"european_aqi_sulphur_dioxide,omitempty"`
	USAQI               float64 `json:"us_aqi,omitempty"`
	USAQIPM25           float64 `json:"us_aqi_pm2_5,omitempty"`
	USAQIPM10           float64 `json:"us_aqi_pm10,omitempty"`
	USAQINO2            float64 `json:"us_aqi_nitrogen_dioxide,omitempty"`
	USAQIOzone          float64 `json:"us_aqi_ozone,omitempty"`
	USAQISO2            float64 `json:"us_aqi_sulphur_dioxide,omitempty"`
	USAQICarbonMonoxide float64 `json:"us_aqi_carbon_monoxide,omitempty"`
}

type CurrentUnits struct {
	PM10                string `json:"pm10,omitempty"`
	PM25                string `json:"pm2_5,omitempty"`
	CarbonMonoxide      string `json:"carbon_monoxide,omitempty"`
	CarbonDioxide       string `json:"carbon_dioxide,omitempty"`
	NitrogenDioxide     string `json:"nitrogen_dioxide,omitempty"`
	SulphurDioxide      string `json:"sulphur_dioxide,omitempty"`
	Ozone               string `json:"ozone,omitempty"`
	AerosolOpticalDepth string `json:"aerosol_optical_depth,omitempty"`
	Dust                string `json:"dust,omitempty"`
	UVIndex             string `json:"uv_index,omitempty"`
	UVIndexClearSky     string `json:"uv_index_clear_sky,omitempty"`
	Ammonia             string `json:"ammonia,omitempty"`
	Methane             string `json:"methane,omitempty"`
	AlderPollen         string `json:"alder_pollen,omitempty"`
	BirchPollen         string `json:"birch_pollen,omitempty"`
	GrassPollen         string `json:"grass_pollen,omitempty"`
	MugwortPollen       string `json:"mugwort_pollen,omitempty"`
	OlivePollen         string `json:"olive_pollen,omitempty"`
	RagweedPollen       string `json:"ragweed_pollen,omitempty"`
	EuropeanAQI         string `json:"european_aqi,omitempty"`
	EuropeanAQIPM25     string `json:"european_aqi_pm2_5,omitempty"`
	EuropeanAQIPM10     string `json:"european_aqi_pm10,omitempty"`
	EuropeanAQINO2      string `json:"european_aqi_nitrogen_dioxide,omitempty"`
	EuropeanAQIOzone    string `json:"european_aqi_ozone,omitempty"`
	EuropeanAQISO2      string `json:"european_aqi_sulphur_dioxide,omitempty"`
	USAQI               string `json:"us_aqi,omitempty"`
	USAQIPM25           string `json:"us_aqi_pm2_5,omitempty"`
	USAQIPM10           string `json:"us_aqi_pm10,omitempty"`
	USAQINO2            string `json:"us_aqi_nitrogen_dioxide,omitempty"`
	USAQIOzone          string `json:"us_aqi_ozone,omitempty"`
	USAQISO2            string `json:"us_aqi_sulphur_dioxide,omitempty"`
	USAQICarbonMonoxide string `json:"us_aqi_carbon_monoxide,omitempty"`
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
