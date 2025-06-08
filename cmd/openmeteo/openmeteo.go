package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/meteo"
	"github.com/ninedraft/daily-bacon/internal/models"
)

func main() {
	exitCode := 0
	defer func() {
		if exitCode != 0 {
			_ = os.Stdout.Sync()
			os.Exit(exitCode)
		}
	}()

	log.SetFlags(log.Lshortfile)

	params := meteo.Params{
		Latitude:     34.707130,
		Longitude:    33.022617,
		PastDays:     1,
		ForecastDays: 1,
		Current:      []string{"pm10", "pm2_5", "dust", "european_aqi"},
		Timezone:     "GMT",
	}
	bindRequestFlags(flag.CommandLine, "api", &params)

	timeout := 10 * time.Second
	flag.DurationVar(&timeout, "client.timeout", timeout, "HTTP client timeout")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cl := client.New(http.DefaultClient.Transport)
	meteoClient := meteo.New(cl)

	resp, err := meteoClient.AirQuality(ctx, params)
	if err != nil {
		log.Printf("doing request: %v", err)
		exitCode = 10
		return
	}

	if err := formatData(os.Stdout, resp); err != nil {
		log.Printf("formatting response: %v", err)
		exitCode = 11
	}
	_, _ = os.Stdout.WriteString("\n")
}

func bindRequestFlags(flags *flag.FlagSet, prefix string, p *meteo.Params) {
	if prefix != "" {
		prefix += "."
	}

	flags.Float64Var(&p.Latitude, prefix+"latitude", p.Latitude, "latitude for air quality request")
	flags.Float64Var(&p.Longitude, prefix+"longitude", p.Longitude, "longitude for air quality request")

	flags.Func(prefix+"current", "current variables", func(value string) error {
		fields := strings.FieldsFuncSeq(value, flagSliceField)
		for field := range fields {
			p.Current = append(p.Current, field)
		}
		return nil
	})

	flags.Func(prefix+"hourly", "hourly variables", func(value string) error {
		fields := strings.FieldsFuncSeq(value, flagSliceField)
		for field := range fields {
			p.Hourly = append(p.Hourly, field)
		}
		return nil
	})

	flags.Func(prefix+"daily", "daily variables", func(value string) error {
		fields := strings.FieldsFuncSeq(value, flagSliceField)
		for field := range fields {
			p.Daily = append(p.Daily, field)
		}
		return nil
	})

	flags.Func(prefix+"start-date", "start date (YYYY-MM-DD)", func(value string) error {
		t, err := time.Parse(time.DateOnly, value)
		if err != nil {
			return fmt.Errorf("%q: time.Parse: %w", value, err)
		}
		p.StartDate = t
		return nil
	})

	flags.Func(prefix+"end-date", "end date (YYYY-MM-DD)", func(value string) error {
		t, err := time.Parse(time.DateOnly, value)
		if err != nil {
			return fmt.Errorf("%q: time.Parse: %w", value, err)
		}
		p.EndDate = t
		return nil
	})

	flags.StringVar(&p.Timezone, prefix+"timezone", p.Timezone, "timezone for dates returned")
	flags.IntVar(&p.ForecastDays, prefix+"forecast-days", p.ForecastDays, "number of forecast days")
	flags.IntVar(&p.PastDays, prefix+"past-days", p.PastDays, "number of past days")
}

func flagSliceField(ru rune) bool {
	return strings.ContainsRune(",|", ru) || unicode.IsSpace(ru)
}

// formatFloat prints f with minimal precision (e.g. 1.00→"1", 1.20→"1.2").
func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func formatData(dst io.Writer, data models.AirQualityResponse) error {
	// 0 min width, 4-char tabs, 2 spaces padding, pad char=' ', no flags
	wr := tabwriter.NewWriter(dst, 0, 4, 2, ' ', 0)

	// no data?
	if data.Current == nil {
		fmt.Fprintln(dst, "\nno data")
		return nil
	}
	curr := data.Current
	units := data.CurrentUnits

	// header
	fmt.Fprintln(dst, "\n🕒  Current Air Quality")

	// table-driven all fields
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
			// e.g. "🟤 PM₁₀:    19.7 μg/m³"
			fmt.Fprintf(wr, "%s %s:\t%s %s\n",
				f.icon,
				f.label,
				formatFloat(f.value),
				f.unit,
			)
		}
	}

	if err := wr.Flush(); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}
