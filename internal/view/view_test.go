package view

import (
	"bytes"
	"testing"

	"github.com/ninedraft/daily-bacon/internal/models"
	"github.com/stretchr/testify/require"
)

func TestAirQuality(t *testing.T) {
	var b bytes.Buffer
	err := AirQuality(&b, models.AirQualityResponse{
		Current:      &models.CurrentData{PM10: 1},
		CurrentUnits: &models.CurrentUnits{PM10: "ug/m3"},
	})
	require.NoError(t, err)
	require.NotEmpty(t, b.String())
}
