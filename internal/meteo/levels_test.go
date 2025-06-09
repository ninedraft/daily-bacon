package meteo

import "testing"

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelGood, "Good"},
		{LevelWatch, "Beware"},
		{LevelLimitExceeded, "Limit Exceeded"},
		{LevelActNow, "Act Now"},
		{Level(100), "Unknown"},
	}

	for _, tc := range tests {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q; want %q", tc.level, got, tc.want)
		}
	}
}

func TestLevelOf(t *testing.T) {
	tests := []struct {
		key   string
		value float64
		want  Level
	}{
		// PM2.5 bands
		{PM2_5, 10, LevelGood},
		{PM2_5, 15, LevelWatch},
		{PM2_5, 25, LevelWatch},
		{PM2_5, 30, LevelLimitExceeded},
		{PM2_5, 50, LevelLimitExceeded},
		{PM2_5, 60, LevelActNow},

		// UVIndex bands
		{UVIndex, 2, LevelGood},
		{UVIndex, 3, LevelWatch},
		{UVIndex, 7, LevelLimitExceeded},
		{UVIndex, 8, LevelActNow},

		// Unknown key defaults to Good
		{"unknown_key", 1000, LevelGood},
	}

	for _, tc := range tests {
		got := LevelOf(tc.key, tc.value)
		if got != tc.want {
			t.Errorf("LevelOf(%q, %v) = %v; want %v", tc.key, tc.value, got, tc.want)
		}
	}
}
