package aq

import "fmt"

// Report describes air quality statistics.
type Report struct {
	Location string
	AQI      int
}

// Format produces human readable report string.
func Format(r Report) string {
	return fmt.Sprintf("Air quality at %s: AQI %d", r.Location, r.AQI)
}
