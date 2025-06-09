package view

import (
	"fmt"
	"io"
	"strconv"
	"text/tabwriter"

	"github.com/ninedraft/daily-bacon/internal/meteo"
	"github.com/ninedraft/daily-bacon/internal/models"
)

const (
	tabWidth = 4
	tabPad   = 2
)

func AirQuality(dst io.Writer, data models.AirQualityResponse) error {
	wr := tabwriter.NewWriter(dst, 0, tabWidth, tabPad, ' ', 0)

	if data.Current == nil {
		fmt.Fprintln(dst, "no data")
		return nil
	}
	curr := data.Current
	units := data.CurrentUnits

	fmt.Fprintln(dst, "🕒  Current Air Quality")

	type field struct {
		icon, label string
		value       float64
		unit        string
	}
	fields := []field{
		{"🟤", "PM₁₀", curr.PM10, units.PM10},
		{"🔴", "PM₂.₅", curr.PM25, units.PM25},
		{"🛢️", "CO", curr.CarbonMonoxide, units.CarbonMonoxide},
		{"☁️", "CO₂", curr.CarbonDioxide, units.CarbonDioxide},
		{"💨", "NO₂", curr.NitrogenDioxide, units.NitrogenDioxide},
		{"🛑", "SO₂", curr.SulphurDioxide, units.SulphurDioxide},
		{"🟢", "Ozone", curr.Ozone, units.Ozone},
		{"🌫️", "Aerosol Opt. Depth", curr.AerosolOpticalDepth, units.AerosolOpticalDepth},
		{"💨", "Dust", curr.Dust, units.Dust},
		{"🔆", "UV Index", curr.UVIndex, units.UVIndex},
		{"☀️", "UV Index Clear Sky", curr.UVIndexClearSky, units.UVIndexClearSky},
		{"🧪", "Ammonia", curr.Ammonia, units.Ammonia},
		{"🛢️", "Methane", curr.Methane, units.Methane},
		{"🌳", "Alder Pollen", curr.AlderPollen, units.AlderPollen},
		{"🌳", "Birch Pollen", curr.BirchPollen, units.BirchPollen},
		{"🌱", "Grass Pollen", curr.GrassPollen, units.GrassPollen},
		{"🌾", "Mugwort Pollen", curr.MugwortPollen, units.MugwortPollen},
		{"🫒", "Olive Pollen", curr.OlivePollen, units.OlivePollen},
		{"🍂", "Ragweed Pollen", curr.RagweedPollen, units.RagweedPollen},
		{"📊", "EU AQI", curr.EuropeanAQI, units.EuropeanAQI},
		{"📊", "EU AQI PM₂.₅", curr.EuropeanAQIPM25, units.EuropeanAQIPM25},
		{"📊", "EU AQI PM₁₀", curr.EuropeanAQIPM10, units.EuropeanAQIPM10},
		{"📊", "EU AQI NO₂", curr.EuropeanAQINO2, units.EuropeanAQINO2},
		{"📊", "EU AQI Ozone", curr.EuropeanAQIOzone, units.EuropeanAQIOzone},
		{"📊", "EU AQI SO₂", curr.EuropeanAQISO2, units.EuropeanAQISO2},
		{"📊", "US AQI", curr.USAQI, units.USAQI},
		{"📊", "US AQI PM₂.₅", curr.USAQIPM25, units.USAQIPM25},
		{"📊", "US AQI PM₁₀", curr.USAQIPM10, units.USAQIPM10},
		{"📊", "US AQI NO₂", curr.USAQINO2, units.USAQINO2},
		{"📊", "US AQI Ozone", curr.USAQIOzone, units.USAQIOzone},
		{"📊", "US AQI SO₂", curr.USAQISO2, units.USAQISO2},
		{"📊", "US AQI CO", curr.USAQICarbonMonoxide, units.USAQICarbonMonoxide},
	}

	for _, f := range fields {
		if f.value != 0 {
			level := meteo.LevelOf(f.label, f.value)
			levelIcon := "✅"
			switch level {
			case meteo.LevelWatch:
				levelIcon = "😷"
			case meteo.LevelLimitExceeded:
				levelIcon = "⚠️"
			case meteo.LevelActNow:
				levelIcon = "‼️☠️"
			default:
				// pass
			}
			fmt.Fprintf(wr, "%s\t%s:\t%s\t%s\t%s\t%s\n",
				f.icon,
				f.label,
				formatFloat(f.value),
				f.unit,
				level.String(),
				levelIcon,
			)
		}
	}

	if err := wr.Flush(); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
