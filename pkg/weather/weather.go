package weather

import (
	"fmt"
	"time"

	"driffaud.fr/odin/pkg/util"
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
	if len(m.weatherData.Hourly.Time) == 0 {
		return util.BorderStyle.
			Width(m.width-2).
			Height(m.height-2).
			Padding(0, 2).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No weather data available")
	}

	forecastData := GenerateForecastData(m.weatherData)

	favoriteStatus := ""
	if m.isFavorite {
		favoriteStatus = "⭐️"
	} else {
		favoriteStatus = "❌"
	}

	title := fmt.Sprintf("Météo à %s %s", m.placeName, favoriteStatus)

	astroSection := formatAstroInfo(forecastData, m.weatherData.Latitude, m.weatherData.Longitude)
	var forecastSection string
	if len(m.weatherData.Hourly.Time) >= 24 {
		forecastSection = formatForecast(forecastData)
	}

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

func formatAstroInfo(forecastData []ForecastHour, lat, lon float64) string {
	sunInfo := GetSunInfo(lat, lon)
	moonInfo := GetMoonInfo(lat, lon)
	nightForecast := AnalyzeNightForecast(forecastData, sunInfo.Sunset, sunInfo.Sunrise)

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

	precipAndSeeing := fmt.Sprintf("Risque de précipitation: %d%% | Indice de seeing: %d/5",
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

func formatForecast(forecastData []ForecastHour) string {
	title := util.SubtitleStyle.Render("Prévisions des prochaines heures:")

	now := time.Now()
	startIndex := 0

	for i, hour := range forecastData {
		if hour.DateTime.After(now) {
			startIndex = i
			break
		}
	}

	hoursToShow := 24
	if startIndex+hoursToShow > len(forecastData) {
		hoursToShow = len(forecastData) - startIndex
	}

	columns := []table.Column{
		{Title: "Heure", Width: 7},
		{Title: "Nuages", Width: 7},
		{Title: "Pluie", Width: 7},
		{Title: "Seeing", Width: 7},
		{Title: "Vent", Width: 9},
		{Title: "Humidité", Width: 9},
		{Title: "Temp", Width: 7},
		{Title: "Rosée", Width: 7},
	}

	var rows []table.Row
	for i := range hoursToShow {
		idx := startIndex + i
		hour := forecastData[idx]

		row := table.Row{
			hour.DateTime.Format("15h"),
			fmt.Sprintf("%d%%", hour.Clouds),
			fmt.Sprintf("%d%%", hour.PrecipitationProbability),
			fmt.Sprintf("%d/5", hour.Seeing),
			fmt.Sprintf("%.1f km/h", hour.WindSpeed),
			fmt.Sprintf("%d%%", hour.Humidity),
			fmt.Sprintf("%.1f°C", hour.Temperature),
			fmt.Sprintf("%.1f°C", hour.DewPoint),
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
