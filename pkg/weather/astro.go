package weather

import (
	"time"

	"github.com/sixdouglas/suncalc"
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
	Moonrise     time.Time
	Moonset      time.Time
}

// MoonPhaseInfo contains the name and emoji for a moon phase
type moonPhaseInfo struct {
	name  string
	emoji string
}

// GetSunInfo calculates sun-related astronomical information
func GetSunInfo(lat, lon float64) SunInfo {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	observer := suncalc.Observer{Latitude: lat, Longitude: lon, Location: time.Local}
	sunTimes := suncalc.GetTimesWithObserver(today, observer)
	sunTimesTomorrow := suncalc.GetTimesWithObserver(tomorrow, observer)

	return SunInfo{
		Sunset:  sunTimes[suncalc.Sunset].Value,
		Dusk:    sunTimes[suncalc.NauticalDusk].Value,
		Dawn:    sunTimesTomorrow[suncalc.NauticalDawn].Value,
		Sunrise: sunTimesTomorrow[suncalc.Sunrise].Value,
	}
}

// GetMoonInfo calculates moon-related astronomical information
func GetMoonInfo(lat, lon float64) MoonInfo {
	today := time.Now()
	// tomorrow := today.AddDate(0, 0, 1)

	moonTimes := suncalc.GetMoonTimes(today, lat, lon, false)
	// moonTimesTomorrow := suncalc.GetMoonTimesWithObserver(tomorrow, observer)
	phase := suncalc.GetMoonIllumination(today)
	phaseInfo := getMoonPhaseInfo(phase)
	illumination := phase.Fraction * 100

	return MoonInfo{
		PhaseName:    phaseInfo.name,
		PhaseEmoji:   phaseInfo.emoji,
		Illumination: illumination,
		Moonrise:     moonTimes.Rise,
		Moonset:      moonTimes.Set,
	}
}

// getMoonPhaseInfo determines the moon phase name and emoji
func getMoonPhaseInfo(phase suncalc.MoonIllumination) moonPhaseInfo {
	switch {
	case phase.Phase < 0.125 || phase.Phase >= 0.875: // New moon
		return moonPhaseInfo{"Nouvelle lune", "🌑"}
	case phase.Phase < 0.25: // Waxing crescent
		return moonPhaseInfo{"Premier croissant", "🌒"}
	case phase.Phase < 0.375: // First quarter
		return moonPhaseInfo{"Premier quartier", "🌓"}
	case phase.Phase < 0.5: // Waxing gibbous
		return moonPhaseInfo{"Gibbeuse croissante", "🌔"}
	case phase.Phase < 0.625: // Full moon
		return moonPhaseInfo{"Pleine lune", "🌕"}
	case phase.Phase < 0.75: // Waning gibbous
		return moonPhaseInfo{"Gibbeuse décroissante", "🌖"}
	case phase.Phase < 0.875: // Last quarter
		return moonPhaseInfo{"Dernier quartier", "🌗"}
	default: // Waning crescent
		return moonPhaseInfo{"Dernier croissant", "🌘"}
	}
}
