package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
)

const openMeteoAPI = "https://api.open-meteo.com/v1/forecast"

// WeatherData represents the entire weather data response
type WeatherData struct {
	Current        CurrentWeather `json:"current"`
	CurrentUnits   CurrentUnits   `json:"current_units"`
	Hourly         HourlyWeather  `json:"hourly"`
	HourlyUnits    HourlyUnits    `json:"hourly_units"`
	Daily          DailyWeather   `json:"daily"`
	DailyUnits     DailyUnits     `json:"daily_units"`
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	Elevation      float64        `json:"elevation"`
	GenerationTime float64        `json:"generationtime_ms"`
	Timezone       string         `json:"timezone"`
	TimezoneAbbr   string         `json:"timezone_abbreviation"`
}

type CurrentWeather struct {
	Time                     string  `json:"time"`
	Temperature              float64 `json:"temperature_2m"`
	RelativeHumidity         int     `json:"relative_humidity_2m"`
	CloudCover               int     `json:"cloud_cover"`
	WindSpeed                float64 `json:"wind_speed_10m"`
	WindDirection            float64 `json:"wind_direction_10m"`
	PrecipitationProbability int     `json:"precipitation_probability"`
	DewPoint                 float64 `json:"dew_point_2m"`
}

type CurrentUnits struct {
	Temperature              string `json:"temperature_2m"`
	RelativeHumidity         string `json:"relative_humidity_2m"`
	CloudCover               string `json:"cloud_cover"`
	WindSpeed                string `json:"wind_speed_10m"`
	WindDirection            string `json:"wind_direction_10m"`
	PrecipitationProbability string `json:"precipitation_probability"`
	DewPoint                 string `json:"dew_point_2m"`
}

type HourlyWeather struct {
	Time                     []string  `json:"time"`
	Temperature              []float64 `json:"temperature_2m"`
	RelativeHumidity         []int     `json:"relative_humidity_2m"`
	CloudCover               []int     `json:"cloud_cover"`
	CloudCoverLow            []int     `json:"cloud_cover_low"`
	CloudCoverMid            []int     `json:"cloud_cover_mid"`
	CloudCoverHigh           []int     `json:"cloud_cover_high"`
	WindSpeed                []float64 `json:"wind_speed_10m"`
	WindDirection            []float64 `json:"wind_direction_10m"`
	PrecipitationProbability []int     `json:"precipitation_probability"`
	DewPoint                 []float64 `json:"dew_point_2m"`
}

type HourlyUnits struct {
	Temperature              string `json:"temperature_2m"`
	RelativeHumidity         string `json:"relative_humidity_2m"`
	CloudCover               string `json:"cloud_cover"`
	CloudCoverLow            string `json:"cloud_cover_low"`
	CloudCoverMid            string `json:"cloud_cover_mid"`
	CloudCoverHigh           string `json:"cloud_cover_high"`
	WindSpeed                string `json:"wind_speed_10m"`
	WindDirection            string `json:"wind_direction_10m"`
	PrecipitationProbability string `json:"precipitation_probability"`
	DewPoint                 string `json:"dew_point_2m"`
}

type DailyWeather struct {
	Time    []string `json:"time"`
	Sunrise []string `json:"sunrise"`
	Sunset  []string `json:"sunset"`
}

type DailyUnits struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
}

// WeatherResultMsg holds weather data for the update function
type WeatherResultMsg struct {
	Data WeatherData
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

		var weather WeatherData
		if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode weather data: %w", err))
		}

		return WeatherResultMsg{Data: weather}
	}
}
