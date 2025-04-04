package photon

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"driffaud.fr/odin/internal/domain"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const photonAPI = "https://photon.komoot.io/api"

// PhotonResponse represents the API response from Photon
type PhotonResponse struct {
	Features []struct {
		Properties struct {
			Name     string `json:"name"`
			City     string `json:"city,omitempty"`
			State    string `json:"state,omitempty"`
			Country  string `json:"country,omitempty"`
			Street   string `json:"street,omitempty"`
			PostCode string `json:"postCode,omitempty"`
		} `json:"properties"`
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

// SearchResultsMsg carries search results back to the model
type SearchResultsMsg []list.Item

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

			items = append(items, domain.Place{
				Name:      name,
				Address:   address,
				Latitude:  lat,
				Longitude: lon,
			})
		}

		if len(items) == 0 {
			return ErrMsg(fmt.Errorf("no results found for %s", query))
		}

		return SearchResultsMsg(items)
	}
}
