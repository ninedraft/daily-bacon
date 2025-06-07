package meteo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ninedraft/daily-bacon/internal/client"
	"github.com/ninedraft/daily-bacon/internal/models"
	"github.com/stretchr/testify/require"
)

func TestClient_AirQuality(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/", r.URL.Path)
		require.Equal(t, "34.700000", r.URL.Query().Get("latitude"))
		require.Equal(t, "33.020000", r.URL.Query().Get("longitude"))
		require.Equal(t, "pm2_5,ozone", r.URL.Query().Get("hourly"))
		_ = json.NewEncoder(w).Encode(models.AirQualityResponse{Latitude: 34.7, Longitude: 33.02})
	}))
	defer srv.Close()

	c := New(client.New(srv.Client().Transport))
	c.url = srv.URL

	resp, err := c.AirQuality(context.Background(), Params{
		Latitude:  34.7,
		Longitude: 33.02,
		Hourly:    "pm2_5,ozone",
	})
	require.NoError(t, err)
	require.Equal(t, 34.7, resp.Latitude)
	require.Equal(t, 33.02, resp.Longitude)
}
