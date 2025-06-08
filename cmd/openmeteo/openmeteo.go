package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/meteo"
	"github.com/ninedraft/daily-bacon/internal/view"
)

const (
	defaultLatitude  = 34.707130
	defaultLongitude = 33.022617
	defaultTimeout   = 10 * time.Second
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
		Latitude:     defaultLatitude,
		Longitude:    defaultLongitude,
		PastDays:     1,
		ForecastDays: 1,
		Current:      []string{"pm10", "pm2_5", "dust", "european_aqi"},
		Timezone:     "GMT",
	}
	bindRequestFlags(flag.CommandLine, "api", &params)

	timeout := defaultTimeout
	flag.DurationVar(&timeout, "client.timeout", timeout, "HTTP client timeout")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cl := client.New(http.DefaultClient.Transport)
	meteoClient := meteo.New(cl)

	resp, err := meteoClient.AirQuality(ctx, params)
	if err != nil {
		log.Printf("doing request: %v", err)
		exitCode = 10
		return
	}

	if err := view.AirQuality(os.Stdout, resp); err != nil {
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
