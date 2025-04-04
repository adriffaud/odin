package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"driffaud.fr/odin/internal/domain"
)

const openMeteoAPI = "https://api.open-meteo.com/v1/forecast"

func GetWeather(lat, lon float64) (domain.WeatherData, error) {
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

	var weather domain.WeatherData

	resp, err := http.Get(url)
	if err != nil {
		return weather, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return weather, fmt.Errorf("openmeteo API returned non-200 status: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return weather, fmt.Errorf("failed to decode weather data: %w", err)
	}

	return weather, nil
}
