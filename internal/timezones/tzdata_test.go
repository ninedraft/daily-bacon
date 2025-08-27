package timezones

import (
	"testing"
)

func TestAllLocationS(t *testing.T) {
	t.Parallel()

	if len(allTimezones()) == 0 {
		t.Fatal("no timezones found")
	}
}
