package ui

import (
	"fmt"
	"time"

	"driffaud.fr/odin/internal/domain"
	"driffaud.fr/odin/internal/domain/astro"
	"driffaud.fr/odin/internal/forecast"
	"driffaud.fr/odin/internal/i18n"
	"driffaud.fr/odin/internal/platform/storage"
	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WeatherModel represents the weather view component
type WeatherModel struct {
	width, height int
	weatherData   domain.WeatherData
	placeName     string
	isFavorite    bool
	favorites     *storage.FavoritesStore
	selectedPlace domain.Place
}

// NewWeatherModel creates a new weather view model
func NewWeatherModel(data domain.WeatherData, place domain.Place, favorites *storage.FavoritesStore, width, height int) WeatherModel {
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
			Padding(0, 1).
			Align(lipgloss.Center, lipgloss.Center).
			Render(i18n.T("weather.no_data", nil))
	}

	forecastData := forecast.GenerateForecastData(m.weatherData)

	astroSection := formatAstroInfo(forecastData, m.weatherData.Latitude, m.weatherData.Longitude)
	var forecastSection string
	if len(m.weatherData.Hourly.Time) >= 24 {
		forecastSection = formatForecast(forecastData)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.headerView(),
		astroSection,
		forecastSection,
		helpView,
	)

	return util.BorderStyle.
		Width(m.width-2).
		Height(m.height-2).
		Padding(0, 1).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func (m WeatherModel) headerView() string {
	favoriteStatus := ""
	if m.isFavorite {
		favoriteStatus = "⭐️"
	} else {
		favoriteStatus = "❌"
	}

	title := fmt.Sprint(i18n.T("weather.at_title", map[string]any{
		"Place":     m.placeName,
		"FavStatus": favoriteStatus,
	}))

	return util.TitleStyle.Render(title)
}

func formatTime(t time.Time) string {
	return t.Format("15:04")
}

func formatAstroInfo(forecastData []forecast.ForecastHour, lat, lon float64) string {
	sunInfo := astro.GetSunInfo(lat, lon)
	moonInfo := astro.GetMoonInfo(lat, lon)
	nightForecast := forecast.AnalyzeNightForecast(forecastData, sunInfo.Sunset, sunInfo.Sunrise)

	sunInfoStr := fmt.Sprint(i18n.T("weather.sunset", map[string]any{
		"Sunset":  formatTime(sunInfo.Sunset),
		"Dusk":    formatTime(sunInfo.Dusk),
		"Dawn":    formatTime(sunInfo.Dawn),
		"Sunrise": formatTime(sunInfo.Sunrise),
	}),
	)

	moonInfoStr := fmt.Sprint(i18n.T("weather.moonphase", map[string]any{
		"MoonEmoji":    moonInfo.PhaseEmoji,
		"Moonrise":     formatTime(moonInfo.Moonrise),
		"Moonset":      formatTime(moonInfo.Moonset),
		"Illumination": fmt.Sprintf("%.0f", moonInfo.Illumination),
		"PhaseName":    moonInfo.PhaseName,
	}),
	)

	forecastTitle := i18n.T("weather.conditions_title", nil)

	var observationTimeStr string
	if nightForecast.BestObservation.TimeRange != nil {
		observationTimeStr = fmt.Sprint(i18n.T("weather.best_period", map[string]any{
			"Start":      nightForecast.BestObservation.TimeRange.Start,
			"End":        nightForecast.BestObservation.TimeRange.End,
			"CloudCover": nightForecast.BestObservation.LowestCloudCover,
		}))
	} else {
		observationTimeStr = fmt.Sprint(i18n.T("weather.unfavorable", map[string]any{
			"CloudCover": nightForecast.DisplayCloudCover,
		}))
	}

	weatherConditions := fmt.Sprint(i18n.T("weather.conditions", map[string]any{
		"Temp":      nightForecast.NightlyTemperature,
		"Humidity":  nightForecast.NightlyHumidity,
		"WindSpeed": nightForecast.NightlyWindSpeed,
		"WindDir":   nightForecast.WindDirectionText,
		"DewPoint":  nightForecast.NightlyDewPoint,
	}))

	precipAndSeeing := fmt.Sprint(i18n.T("weather.precip_and_seeing", map[string]any{
		"Precip": nightForecast.MaxPrecipProbability,
		"Seeing": nightForecast.SeeingIndex,
	}))

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

func formatForecast(forecastData []forecast.ForecastHour) string {
	title := util.SubtitleStyle.Render(i18n.T("forecast.title", nil))

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
		{Title: i18n.T("forecast.hour", nil), Width: 7},
		{Title: i18n.T("forecast.clouds", nil), Width: 7},
		{Title: i18n.T("forecast.rain", nil), Width: 7},
		{Title: i18n.T("forecast.seeing", nil), Width: 7},
		{Title: i18n.T("forecast.wind", nil), Width: 9},
		{Title: i18n.T("forecast.humidity", nil), Width: 9},
		{Title: i18n.T("forecast.temp", nil), Width: 7},
		{Title: i18n.T("forecast.dew", nil), Width: 7},
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
