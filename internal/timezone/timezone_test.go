package timezone_test

import (
	"testing"

	"github.com/ninedraft/daily-bacon/internal/timezone"
)

func TestFind(t *testing.T) {
	t.Parallel()

	got := timezone.Find("Nicosia")
	if len(got) == 0 {
		t.Fatal("no results")
	}
	t.Log("got:", got)
}
