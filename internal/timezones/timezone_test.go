package timezones_test

import (
	"testing"

	"github.com/ninedraft/daily-bacon/internal/timezones"
)

func TestFind(t *testing.T) {
	t.Parallel()

	got := timezones.Find("Nicosia")
	if len(got) == 0 {
		t.Fatal("no results")
	}
	t.Log("got:", got)
}
