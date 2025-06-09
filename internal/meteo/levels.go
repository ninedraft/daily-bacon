package meteo

// Level indicates a health‐risk band.
type Level int

const (
	// LevelGood indicates Low / “Good”
	LevelGood Level = iota
	// LevelWatch indicates Medium / “Watch”
	LevelWatch
	// LevelLimitExceeded indicates High / “Limit Exceeded”
	LevelLimitExceeded
	// LevelActNow indicates Dangerous / “Act Now”
	LevelActNow
)

// String returns the textual representation of the Level.
func (l Level) String() string {
	switch l {
	case LevelGood:
		return "Good"
	case LevelWatch:
		return "Watch"
	case LevelLimitExceeded:
		return "Limit Exceeded"
	case LevelActNow:
		return "Act Now"
	default:
		return "Unknown"
	}
}

// registry maps variable keys to three ascending thresholds.
// value ≤ thresholds[0] → Good
// value ≤ thresholds[1] → Watch
// value ≤ thresholds[2] → Limit Exceeded
// value > thresholds[2] → Act Now
var registry = map[string][3]float64{
	PM2_5:               {10, 25, 50},
	PM10:                {20, 50, 100},
	NitrogenDioxide:     {40, 120, 230},
	SulphurDioxide:      {100, 350, 500},
	Ozone:               {50, 130, 240},
	CarbonMonoxide:      {4, 10, 35},
	CarbonDioxide:       {1000, 2000, 5000},
	Ammonia:             {100, 200, 600},
	AerosolOpticalDepth: {0.10, 0.30, 0.80},
	Methane:             {5, 50, 50000},
	UVIndex:             {2, 5, 7},
	Dust:                {20, 50, 100},
	AlderPollen:         {20, 50, 100},
	BirchPollen:         {20, 50, 100},
	GrassPollen:         {20, 50, 100},
	MugwortPollen:       {20, 50, 100},
	OlivePollen:         {20, 50, 100},
	RagweedPollen:       {20, 50, 100},
}

// Level returns the health‐risk band for the given key and value.
// If the key is not registered, it returns LevelGood.
func LevelOf(key string, value float64) Level {
	if b, ok := registry[key]; ok {
		switch {
		case value <= b[0]:
			return LevelGood
		case value <= b[1]:
			return LevelWatch
		case value <= b[2]:
			return LevelLimitExceeded
		default:
			return LevelActNow
		}
	}
	return LevelGood
}
