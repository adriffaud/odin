package app

import (
	"fmt"
	"time"

	"driffaud.fr/odin/pkg/util"
	"driffaud.fr/odin/pkg/weather"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// InitInput initializes the text input component
func InitInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Entrer un nom de lieu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	return ti
}

// InitFavoritesList initializes the favorites list component
func InitFavoritesList(items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Lieux favoris"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = util.TitleStyle
	l.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)
	return l
}

// InitResults initializes the results list component
func InitResultsList() list.Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Résultats"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = util.TitleStyle
	l.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)
	return l
}

// RenderPlaces renders the input screen
func RenderPlaces(input textinput.Model, favoritesList list.Model, focusIndex int, width, height int) string {
	title := util.TitleStyle.Render("Météo astronomique")

	inputTitle := "Rechercher un lieu"
	if focusIndex == 0 {
		inputTitle = "> " + inputTitle + " <"
	}
	inputTitleStyled := lipgloss.NewStyle().Bold(true).Render(inputTitle)
	inputField := lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Render(input.View())

	var favoritesSection string

	if len(favoritesList.Items()) > 0 {
		favoritesTitle := "Favoris"
		if focusIndex == 1 {
			favoritesTitle = "> " + favoritesTitle + " <"
		}

		favoritesSection = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render(favoritesTitle),
			favoritesList.View(),
		)
	} else {
		favoritesSection = lipgloss.NewStyle().
			Faint(true).
			Render("Aucun lieu favori - Appuyez sur F2 pour en ajouter")
	}

	helpText := lipgloss.NewStyle().
		Faint(true).
		Render("Tab : changer de focus | Entrée : sélectionner | Esc : quitter")

	inputSection := lipgloss.JoinVertical(lipgloss.Left,
		inputTitleStyled,
		inputField)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, inputSection, favoritesSection),
		"",
		helpText,
	)

	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// RenderResults renders the results screen
func RenderResults(resultsList list.Model, width, height int) string {
	hint := "(entrer pour sélectionner, esc pour quitter)"
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		resultsList.View(),
		"",
		hint,
	)

	return util.BorderStyle.
		Width(width - 2).
		Height(height - 2).
		Render(content)
}

// RenderLoading renders the loading screen
func RenderLoading(spinnerView string, width, height int) string {
	loadingMessage := fmt.Sprintf("%s Chargement...", spinnerView)
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(loadingMessage)
}

// RenderError renders the error screen
func RenderError(err error, width, height int) string {
	errorMsg := fmt.Sprintf("Erreur: %s\n\nAppuyer sur une touche pour continuer...", err)
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(errorMsg)
}

// RenderWeather renders the weather screen
func RenderWeather(weather weather.WeatherData, placeName string, isFavorite bool, width, height int) string {
	forecastSection := formatForecast(weather)
	astroSection := formatAstroInfo(weather)

	favoriteStatus := ""
	if isFavorite {
		favoriteStatus = "⭐️"
	} else {
		favoriteStatus = "❌"
	}

	title := fmt.Sprintf("Météo à %s %s", placeName, favoriteStatus)

	helpText := "ESC : retourner au menu principal"
	if isFavorite {
		helpText = "F3 : retirer des favoris | " + helpText
	} else {
		helpText = "F2 : ajouter aux favoris | " + helpText
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		util.TitleStyle.Render(title),
		astroSection,
		forecastSection,
		helpText,
	)

	return util.BorderStyle.
		Width(width - 2).
		Height(height - 2).
		PaddingLeft(2).
		Render(content)
}

func formatAstroInfo(w weather.WeatherData) string {
	if len(w.Hourly.Time) == 0 {
		return ""
	}

	lat := w.Latitude
	lon := w.Longitude

	sunInfo := weather.GetSunInfo(lat, lon)
	moonInfo := weather.GetMoonInfo(lat, lon)
	nightForecast := weather.AnalyzeNightForecast(w, sunInfo.Sunset, sunInfo.Sunrise)

	sunInfoStr := fmt.Sprintf("☀️ Coucher : %s | Crépuscule astro : %s | Aube astro : %s | Lever : %s",
		formatTime(sunInfo.Sunset),
		formatTime(sunInfo.Dusk),
		formatTime(sunInfo.Dawn),
		formatTime(sunInfo.Sunrise),
	)

	moonInfoStr := fmt.Sprintf("%s Lever : %s | Coucher: %s | Illumination : %.0f%% (%s)",
		moonInfo.PhaseEmoji,
		formatTime(moonInfo.Moonrise),
		formatTime(moonInfo.Moonset),
		moonInfo.Illumination,
		moonInfo.PhaseName,
	)

	forecastTitle := "🔭 Conditions d'observation cette nuit:"

	var observationTimeStr string
	if nightForecast.BestObservation.TimeRange != nil {
		observationTimeStr = fmt.Sprintf("Meilleure période: %dh à %dh (couverture nuageuse: %d%%)",
			nightForecast.BestObservation.TimeRange.Start,
			nightForecast.BestObservation.TimeRange.End,
			nightForecast.BestObservation.LowestCloudCover)
	} else {
		observationTimeStr = fmt.Sprintf("Conditions défavorables (couverture nuageuse: %d%%)",
			nightForecast.DisplayCloudCover)
	}

	weatherConditions := fmt.Sprintf("Temp: %d°C | Humidité: %d%% | Vent: %d km/h %s | Point de rosée: %d°C",
		nightForecast.NightlyTemperature,
		nightForecast.NightlyHumidity,
		nightForecast.NightlyWindSpeed,
		nightForecast.WindDirectionText,
		nightForecast.NightlyDewPoint)

	precipAndSeeing := fmt.Sprintf("Risque de précipitation: %d%% | Indice de seeing: %d/20",
		nightForecast.MaxPrecipProbability,
		nightForecast.SeeingIndex)

	nightForecastStr := lipgloss.JoinVertical(
		lipgloss.Left,
		forecastTitle,
		observationTimeStr,
		weatherConditions,
		precipAndSeeing,
	)

	return util.AstroInfoStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		sunInfoStr,
		"",
		moonInfoStr,
		"",
		nightForecastStr,
	))
}

func formatTime(t time.Time) string {
	return t.Format("15:04")
}

func formatForecast(w weather.WeatherData) string {
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
		t, _ := time.Parse(util.ISO8601Format, timeStr)
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
		t, _ := time.Parse(util.ISO8601Format, timeStr)

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

	tableView := util.TableStyle.Render(t.View())

	return util.WeatherSectionStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		util.WeatherInfoStyle.Render(tableView),
	))
}
