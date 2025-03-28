package types

import "github.com/charmbracelet/bubbles/list"

// Place represents a location with a name and address
type Place struct {
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
}

// Implement list.Item interface
func (p Place) Title() string       { return p.Name }
func (p Place) Description() string { return p.Address }
func (p Place) FilterValue() string { return p.Name }

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
