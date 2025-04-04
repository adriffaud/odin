package domain

// Place represents a location with a name, address and coordinates
type Place struct {
	Name      string
	Address   string
	Latitude  float64
	Longitude float64
}

func (p Place) Title() string       { return p.Name }
func (p Place) Description() string { return p.Address }
func (p Place) FilterValue() string { return p.Name }
