package weather

import (
	"fmt"
	"time"

	"driffaud.fr/odin/pkg/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WeatherModel represents the weather view component
type WeatherModel struct {
	width, height int
	weatherData   WeatherData
	placeName     string
	isFavorite    bool
	favorites     *FavoritesStore
	selectedPlace Place
}

// NewWeatherModel creates a new weather view model
func NewWeatherModel(data WeatherData, place Place, favorites *FavoritesStore, width, height int) WeatherModel {
	isFavorite := favorites.IsFavorite(place)
	placeName := place.Name + " (" + place.Address + ")"
	helpModel := help.New()
	helpModel.ShowAll = true

	return WeatherModel{
		width:         width,
		height:        height,
		weatherData:   data,
		placeName:     placeName,
		isFavorite:    isFavorite,
		favorites:     favorites,
		selectedPlace: place,
	}
}

// Init initializes the weather model
func (m WeatherModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the weather model
func (m WeatherModel) Update(msg tea.Msg) (WeatherModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the weather display
func (m WeatherModel) View(helpView string) string {
	forecastSection := formatForecast(m.weatherData)
	astroSection := formatAstroInfo(m.weatherData)

	favoriteStatus := ""
	if m.isFavorite {
		favoriteStatus = "⭐️"
	} else {
		favoriteStatus = "❌"
	}

	title := fmt.Sprintf("Météo à %s %s", m.placeName, favoriteStatus)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		util.TitleStyle.Render(title),
		astroSection,
		forecastSection,
		helpView,
	)

	return util.BorderStyle.
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// BackToPlaceMsg is used to signal a return to the place view
type BackToPlaceMsg struct{}

// backToPlaceCmd creates a command to go back to the place view
func backToPlaceCmd() tea.Msg {
	return BackToPlaceMsg{}
}

func formatAstroInfo(w WeatherData) string {
	if len(w.Hourly.Time) == 0 {
		return ""
	}

	lat := w.Latitude
	lon := w.Longitude

	sunInfo := GetSunInfo(lat, lon)
	moonInfo := GetMoonInfo(lat, lon)
	nightForecast := AnalyzeNightForecast(w, sunInfo.Sunset, sunInfo.Sunrise)

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

func formatForecast(w WeatherData) string {
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
