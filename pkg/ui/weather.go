package ui

import (
	"fmt"
	"time"

	"driffaud.fr/odin/pkg/types"
	"github.com/charmbracelet/bubbles/table"
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

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))
)

func WeatherView(weather types.WeatherData, placeName string, width, height int) string {
	forecastSection := formatForecast(weather)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render(fmt.Sprintf("Météo à %s", placeName)),
		forecastSection,
	)

	return content
}

func formatForecast(w types.WeatherData) string {
	if len(w.Hourly.Time) < 24 {
		return ""
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Render("Prévisions des prochaines heures:")

	now := time.Now()
	startIndex := 0

	for i, timeStr := range w.Hourly.Time {
		t, _ := time.Parse(TIME_FORMAT, timeStr)
		if t.After(now) {
			startIndex = i
			break
		}
	}

	hoursToShow := 24
	if startIndex+hoursToShow > len(w.Hourly.Time) {
		hoursToShow = len(w.Hourly.Time) - startIndex
	}

	columns := []table.Column{
		{Title: "Heure", Width: 7},
		{Title: "Nuages", Width: 7},
		{Title: "Pluie", Width: 7},
		{Title: "Vent", Width: 9},
		{Title: "Humidité", Width: 9},
		{Title: "Temp", Width: 7},
		{Title: "Rosée", Width: 7},
	}

	var rows []table.Row
	for i := range hoursToShow {
		idx := startIndex + i
		timeStr := w.Hourly.Time[idx]
		t, _ := time.Parse(TIME_FORMAT, timeStr)

		row := table.Row{
			t.Format("15h"),
			fmt.Sprintf("%d%%", w.Hourly.CloudCover[idx]),
			fmt.Sprintf("%d%%", w.Hourly.PrecipitationProbability[idx]),
			fmt.Sprintf("%.1f km/h", w.Hourly.WindSpeed[idx]),
			fmt.Sprintf("%d%%", w.Hourly.RelativeHumidity[idx]),
			fmt.Sprintf("%.1f°C", w.Hourly.Temperature[idx]),
			fmt.Sprintf("%.1f°C", w.Hourly.DewPoint[idx]),
		}
		rows = append(rows, row)
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(hoursToShow+1),
	)

	tableView := tableStyle.Render(t.View())

	return WeatherSectionStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		WeatherInfoStyle.Render(tableView),
	))
}
