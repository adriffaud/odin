package weather

import (
	"math"
	"time"

	"github.com/sj14/astral/pkg/astral"
)

// SunInfo holds astronomical information about the sun
type SunInfo struct {
	Sunset  time.Time
	Dusk    time.Time
	Dawn    time.Time
	Sunrise time.Time
}

// MoonInfo holds astronomical information about the moon
type MoonInfo struct {
	PhaseName    string
	PhaseEmoji   string
	Illumination float64
}

// MoonPhaseInfo contains the name and emoji for a moon phase
type moonPhaseInfo struct {
	name  string
	emoji string
}

// GetSunInfo calculates sun-related astronomical information
func GetSunInfo(lat, lon float64) SunInfo {
	observer := astral.Observer{
		Latitude:  lat,
		Longitude: lon,
	}

	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	dusk, _ := astral.Dusk(observer, today, astral.DepressionAstronomical)
	dawn, _ := astral.Dawn(observer, tomorrow, astral.DepressionAstronomical)
	sunset, _ := astral.Sunset(observer, today)
	sunrise, _ := astral.Sunrise(observer, tomorrow)

	return SunInfo{
		Sunset:  sunset,
		Dusk:    dusk,
		Dawn:    dawn,
		Sunrise: sunrise,
	}
}

// GetMoonInfo calculates moon-related astronomical information
func GetMoonInfo() MoonInfo {
	today := time.Now()
	phase := astral.MoonPhase(today)
	phaseInfo := getMoonPhaseInfo(phase)
	illumination := calculateMoonIllumination(phase)

	return MoonInfo{
		PhaseName:    phaseInfo.name,
		PhaseEmoji:   phaseInfo.emoji,
		Illumination: illumination,
	}
}

// getMoonPhaseInfo determines the moon phase name and emoji
func getMoonPhaseInfo(phase float64) moonPhaseInfo {
	switch {
	case phase < 3.5: // New moon
		return moonPhaseInfo{"Nouvelle lune", "🌑"}
	case phase < 7: // Waxing crescent
		return moonPhaseInfo{"Premier croissant", "🌒"}
	case phase < 10.5: // First quarter
		return moonPhaseInfo{"Premier quartier", "🌓"}
	case phase < 14: // Waxing gibbous
		return moonPhaseInfo{"Gibbeuse croissante", "🌔"}
	case phase < 17.5: // Full moon
		return moonPhaseInfo{"Pleine lune", "🌕"}
	case phase < 21: // Waning gibbous
		return moonPhaseInfo{"Gibbeuse décroissante", "🌖"}
	case phase < 24.5: // Last quarter
		return moonPhaseInfo{"Dernier quartier", "🌗"}
	default: // Waning crescent
		return moonPhaseInfo{"Dernier croissant", "🌘"}
	}
}

// calculateMoonIllumination calculates the percentage of moon illumination
func calculateMoonIllumination(phase float64) float64 {
	normalizedPhase := phase / 28.0

	distanceFromFull := math.Abs(normalizedPhase - 0.5)
	if normalizedPhase > 0.5 {
		distanceFromFull = math.Abs(normalizedPhase - 1.5)
	}

	return 100 * (0.5 - distanceFromFull) * 2
}
