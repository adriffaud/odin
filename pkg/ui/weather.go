package ui

import (
	"fmt"
	"time"

	"driffaud.fr/odin/pkg/types"
	"github.com/charmbracelet/lipgloss"
)

const (
	// ISO8601 format
	TIME_FORMAT = "2006-01-02T15:04"
)

var (
	WeatherInfoStyle = lipgloss.NewStyle().
				MarginLeft(4)

	WeatherSectionStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)
)

func WeatherView(weather types.WeatherData, placeName string, width, height int) string {
	// Format current weather
	currentSection := formatCurrentWeather(weather)

	// Format forecast
	forecastSection := formatForecast(weather)

	// Format sunrise/sunset
	sunSection := formatSunriseSunset(weather)

	// Combine all sections
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render(fmt.Sprintf("Météo à %s", placeName)),
		currentSection,
		sunSection,
		forecastSection,
	)

	return content
}

func formatCurrentWeather(w types.WeatherData) string {
	current := w.Current
	units := w.CurrentUnits

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("105")).
		Render("Conditions actuelles:")

	data := []string{
		fmt.Sprintf("Température: %.1f%s", current.Temperature, units.Temperature),
		fmt.Sprintf("Humidité: %d%s", current.RelativeHumidity, units.RelativeHumidity),
		fmt.Sprintf("Point de rosée: %.1f%s", current.DewPoint, units.DewPoint),
		fmt.Sprintf("Couverture nuageuse: %d%s", current.CloudCover, units.CloudCover),
		fmt.Sprintf("Vent: %.1f%s direction %.0f°", current.WindSpeed, units.WindSpeed, current.WindDirection),
		fmt.Sprintf("Probabilité de précipitation: %d%s", current.PrecipitationProbability, units.PrecipitationProbability),
	}

	dataText := WeatherInfoStyle.Render(lipgloss.JoinVertical(lipgloss.Left, data...))

	return WeatherSectionStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, dataText))
}

func formatSunriseSunset(w types.WeatherData) string {
	if len(w.Daily.Time) == 0 {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("220")).
		Render("Lever et coucher du soleil (aujourd'hui):")

	sunrise, _ := time.Parse(TIME_FORMAT, w.Daily.Sunrise[0])
	sunset, _ := time.Parse(TIME_FORMAT, w.Daily.Sunset[0])

	data := []string{
		fmt.Sprintf("Lever du soleil: %s", sunrise.Format("15:04")),
		fmt.Sprintf("Coucher du soleil: %s", sunset.Format("15:04")),
	}

	dataText := WeatherInfoStyle.Render(lipgloss.JoinVertical(lipgloss.Left, data...))

	return WeatherSectionStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, dataText))
}

func formatForecast(w types.WeatherData) string {
	if len(w.Hourly.Time) < 24 {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render("Prévisions des prochaines heures:")

	var forecast []string

	// Get forecasts for the next 24 hours with 3-hour intervals
	for i := 0; i < 24; i += 3 {
		timeStr := w.Hourly.Time[i]
		t, _ := time.Parse(TIME_FORMAT, timeStr)

		forecast = append(forecast, fmt.Sprintf(
			"%s: %.1f°C, Pluie: %d%%, Vent: %.1f km/h",
			t.Format("15:04"),
			w.Hourly.Temperature[i],
			w.Hourly.PrecipitationProbability[i],
			w.Hourly.WindSpeed[i],
		))
	}

	dataText := WeatherInfoStyle.Render(lipgloss.JoinVertical(lipgloss.Left, forecast...))

	return WeatherSectionStyle.Render(lipgloss.JoinVertical(lipgloss.Left, title, dataText))
}
