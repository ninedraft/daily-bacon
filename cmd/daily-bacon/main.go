package main

import (
	"bytes"
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/meteo"
	"github.com/ninedraft/daily-bacon/internal/tg"
	"github.com/ninedraft/daily-bacon/internal/view"
)

const (
	defaultLatitude  = 34.707130
	defaultLongitude = 33.022617
	defaultTimeout   = 10 * time.Second
)

func main() {
	start := time.Now()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var (
		latitude  = flag.Float64("latitude", defaultLatitude, "air quality latitude")
		longitude = flag.Float64("longitude", defaultLongitude, "air quality longitude")
		timeout   = flag.Duration("timeout", defaultTimeout, "request timeout")
	)

	var groupIDs []string
	flag.Func("group-id", "telegram group id (can be set multiple times, comma/space/| separated)", func(value string) error {
		fields := strings.FieldsFuncSeq(value, flagSliceField)
		for field := range fields {
			if field != "" {
				groupIDs = append(groupIDs, field)
			}
		}
		return nil
	})

	flag.Parse()

	tokenFile := os.Getenv("TELEGRAM_TOKEN_FILE")
	if tokenFile == "" {
		logger.Error("TELEGRAM_TOKEN_FILE is not set")
		return
	}

	tokenBytes, err := os.ReadFile(tokenFile)
	if err != nil {
		logger.Error("read token file", slog.Any("err", err))
		return
	}
	token := strings.TrimSpace(string(tokenBytes))
	if err := os.Setenv("TELEGRAM_TOKEN", token); err != nil {
		logger.Error("set TELEGRAM_TOKEN", slog.Any("err", err))
		return
	}

	params := meteo.Params{
		Latitude:     *latitude,
		Longitude:    *longitude,
		PastDays:     1,
		ForecastDays: 1,
		Current: []string{
			meteo.PM2_5,
			meteo.PM10,
			meteo.Dust,
			meteo.OlivePollen,
			meteo.Ozone,
			meteo.NitrogenDioxide,
			meteo.SulphurDioxide,
			meteo.EuropeanAQI,
		},
		Timezone: "GMT",
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	cl := client.New(http.DefaultClient.Transport)
	meteoClient := meteo.New(cl)

	fetchStart := time.Now()
	resp, err := meteoClient.AirQuality(ctx, params)
	if err != nil {
		logger.Error("fetch air quality", slog.Any("err", err))
		return
	}
	fetchDur := time.Since(fetchStart)

	var buf bytes.Buffer
	if err := view.AirQuality(&buf, resp); err != nil {
		logger.Error("format air quality", slog.Any("err", err))
		return
	}
	msg := buf.String()

	tgClient := tg.New(http.DefaultClient)
	var wg sync.WaitGroup
	wg.Add(len(groupIDs))
	for _, id := range groupIDs {
		go func(id string) {
			defer wg.Done()
			if err := tgClient.SendMessage(ctx, id, msg); err != nil {
				logger.Error("send message", slog.String("chat", id), slog.Any("err", err))
			}
		}(id)
	}
	wg.Wait()

	logger.Info("done", slog.Duration("fetch", fetchDur), slog.Duration("total", time.Since(start)))
}

func flagSliceField(ru rune) bool {
	return strings.ContainsRune(",|", ru) || unicode.IsSpace(ru)
}
