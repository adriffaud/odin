package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"driffaud.fr/odin/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

const openMeteoAPI = "https://api.open-meteo.com/v1/forecast"

// WeatherResultMsg holds weather data for the update function
type WeatherResultMsg struct {
	Data domain.WeatherData
}

// ErrMsg wraps errors for use in tea.Msg
type ErrMsg error

func GetWeather(lat, lon float64) tea.Cmd {
	return func() tea.Msg {
		baseURL, _ := url.Parse(openMeteoAPI)
		params := url.Values{}
		params.Add("latitude", fmt.Sprintf("%f", lat))
		params.Add("longitude", fmt.Sprintf("%f", lon))
		params.Add("current", "temperature_2m,relative_humidity_2m,cloud_cover,wind_speed_10m,wind_direction_10m,precipitation_probability,dew_point_2m")
		params.Add("hourly", "precipitation_probability,dew_point_2m,temperature_2m,relative_humidity_2m,cloud_cover,cloud_cover_low,cloud_cover_mid,cloud_cover_high,wind_speed_10m,wind_direction_10m")
		params.Add("daily", "sunrise,sunset")
		params.Add("timezone", "auto")
		params.Add("forecast_days", "7")
		params.Add("models", "best_match")
		baseURL.RawQuery = params.Encode()
		url := baseURL.String()

		resp, err := http.Get(url)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to fetch weather data: %w", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return ErrMsg(fmt.Errorf("API returned non-200 status: %d", resp.StatusCode))
		}

		var weather domain.WeatherData
		if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode weather data: %w", err))
		}

		return WeatherResultMsg{Data: weather}
	}
}
