package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	photonAPI    = "https://photon.komoot.io/api"
	openMeteoAPI = "https://api.open-meteo.com/v1/forecast"
)

// ErrMsg wraps errors for use in tea.Msg
type ErrMsg error

// SearchPlaces searches for places based on the provided query
func SearchPlaces(query string) tea.Cmd {
	return func() tea.Msg {
		params := url.Values{}
		params.Add("q", query)
		params.Add("lang", "fr")
		reqURL := photonAPI + "?" + params.Encode()

		resp, err := http.Get(reqURL)
		if err != nil {
			return ErrMsg(err)
		}
		defer resp.Body.Close()

		var photonResp PhotonResponse
		if err := json.NewDecoder(resp.Body).Decode(&photonResp); err != nil {
			return ErrMsg(err)
		}

		items := []list.Item{}
		for _, feature := range photonResp.Features {
			props := feature.Properties
			name := props.Name
			if name == "" {
				if props.Street != "" {
					name = props.Street
				} else if props.City != "" {
					name = props.City
				} else {
					continue
				}
			}

			addressParts := []string{}
			if props.Street != "" && props.Street != name {
				addressParts = append(addressParts, props.Street)
			}
			if props.City != "" && props.City != name {
				addressParts = append(addressParts, props.City)
			}
			if props.State != "" && props.State != name {
				addressParts = append(addressParts, props.State)
			}
			if props.Country != "" && props.Country != name {
				addressParts = append(addressParts, props.Country)
			}
			if props.PostCode != "" {
				addressParts = append(addressParts, props.PostCode)
			}

			address := strings.Join(addressParts, ", ")

			var lat, lon float64
			if len(feature.Geometry.Coordinates) >= 2 {
				lon = feature.Geometry.Coordinates[0]
				lat = feature.Geometry.Coordinates[1]
			}

			items = append(items, Place{
				Name:      name,
				Address:   address,
				Latitude:  lat,
				Longitude: lon,
			})
		}

		if len(items) == 0 {
			return ErrMsg(fmt.Errorf("No results found for %s", query))
		}

		return SearchResultsMsg(items)
	}
}

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
