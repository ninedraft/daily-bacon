package meteo

//revive:disable:var-naming
const (
	// PM10 - Particulate matter with diameter ≤ 10 µm close to surface (10 m above ground)
	PM10 = "pm10"
	// PM2_5 - Particulate matter with diameter ≤ 2.5 µm close to surface (10 m above ground)
	PM2_5 = "pm2_5"

	// CarbonMonoxide - Atmospheric CO close to surface (10 m above ground)
	CarbonMonoxide = "carbon_monoxide"
	// NitrogenDioxide - Atmospheric NO₂ close to surface (10 m above ground)
	NitrogenDioxide = "nitrogen_dioxide"
	// SulphurDioxide - Atmospheric SO₂ close to surface (10 m above ground)
	SulphurDioxide = "sulphur_dioxide"
	// Ozone - Atmospheric O₃ close to surface (10 m above ground)
	Ozone = "ozone"

	// CarbonDioxide - CO₂ close to surface (10 m above ground)
	CarbonDioxide = "carbon_dioxide"
	// Ammonia - NH₃ concentration close to surface; only available in Europe
	Ammonia = "ammonia"
	// AerosolOpticalDepth - Aerosol optical depth at 550 nm of the entire atmosphere (haze indicator)
	AerosolOpticalDepth = "aerosol_optical_depth"
	// Methane - CH₄ close to surface (10 m above ground)
	Methane = "methane"
	// Dust - Saharan dust particles close to surface level (10 m above ground)
	Dust = "dust"

	// UVIndex - UV index including cloud effects; see ECMWF UV Index recommendation
	UVIndex = "uv_index"
	// UVIndexClearSky - UV index under clear-sky conditions; see ECMWF UV Index recommendation
	UVIndexClearSky = "uv_index_clear_sky"

	// AlderPollen - Alder pollen concentration in grains/m³; only available in Europe
	AlderPollen = "alder_pollen"
	// BirchPollen - Birch pollen concentration in grains/m³; only available in Europe
	BirchPollen = "birch_pollen"
	// GrassPollen - Grass pollen concentration in grains/m³; only available in Europe
	GrassPollen = "grass_pollen"
	// MugwortPollen - Mugwort pollen concentration in grains/m³; only available in Europe
	MugwortPollen = "mugwort_pollen"
	// OlivePollen - Olive pollen concentration in grains/m³; only available in Europe
	OlivePollen = "olive_pollen"
	// RagweedPollen - Ragweed pollen concentration in grains/m³; only available in Europe
	RagweedPollen = "ragweed_pollen"

	// EuropeanAQI - Consolidated European Air Quality Index (0–20 “good” up to >100 “extremely poor”)
	EuropeanAQI = "european_aqi"
	// EuropeanAQI_PM2_5 - European AQI for PM2.5
	EuropeanAQI_PM2_5 = "european_aqi_pm2_5"
	// EuropeanAQI_PM10 - European AQI for PM10
	EuropeanAQI_PM10 = "european_aqi_pm10"
	// EuropeanAQI_NitrogenDioxide - European AQI for NO₂
	EuropeanAQI_NitrogenDioxide = "european_aqi_nitrogen_dioxide"
	// EuropeanAQI_Ozone - European AQI for O₃
	EuropeanAQI_Ozone = "european_aqi_ozone"
	// EuropeanAQI_SulphurDioxide - European AQI for SO₂
	EuropeanAQI_SulphurDioxide = "european_aqi_sulphur_dioxide"

	// USAQI - Consolidated U.S. Air Quality Index (0–50 “good” up to 301–500 “hazardous”)
	USAQI = "us_aqi"
	// USAQI_PM2_5 - U.S. AQI for PM2.5
	USAQI_PM2_5 = "us_aqi_pm2_5"
	// USAQI_PM10 - U.S. AQI for PM10
	USAQI_PM10 = "us_aqi_pm10"
	// USAQI_NitrogenDioxide - U.S. AQI for NO₂
	USAQI_NitrogenDioxide = "us_aqi_nitrogen_dioxide"
	// USAQI_Ozone - U.S. AQI for O₃
	USAQI_Ozone = "us_aqi_ozone"
	// USAQI_SulphurDioxide - U.S. AQI for SO₂
	USAQI_SulphurDioxide = "us_aqi_sulphur_dioxide"
	// USAQI_CarbonMonoxide - U.S. AQI for CO
	USAQI_CarbonMonoxide = "us_aqi_carbon_monoxide"
)
