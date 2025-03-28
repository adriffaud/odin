package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"driffaud.fr/odin/pkg/types"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const photonAPI = "https://photon.komoot.io/api"

// ErrMsg wraps errors for use in tea.Msg
type ErrMsg error

// SearchPlaces searches for places based on the provided query
func SearchPlaces(query string) tea.Cmd {
	return func() tea.Msg {

		params := url.Values{}
		params.Add("q", query)
		reqURL := photonAPI + "?" + params.Encode()

		resp, err := http.Get(reqURL)
		if err != nil {
			return ErrMsg(err)
		}
		defer resp.Body.Close()

		var photonResp types.PhotonResponse
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

			items = append(items, types.Place{
				Name:      name,
				Address:   address,
				Latitude:  lat,
				Longitude: lon,
			})
		}

		if len(items) == 0 {
			return ErrMsg(fmt.Errorf("No results found for %s", query))
		}

		return types.SearchResultsMsg(items)
	}
}
