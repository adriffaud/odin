package app

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all the keybindings for the application
type KeyMap struct {
	Back           key.Binding
	Tab            key.Binding
	Enter          key.Binding
	Quit           key.Binding
	AddFavorite    key.Binding
	RemoveFavorite key.Binding
	State          ApplicationState
}

// ShortHelp returns keybindings to be shown in the mini help view.
// Will be different based on application state
func (k KeyMap) ShortHelp() []key.Binding {
	switch k.State {
	case StatePlace:
		return []key.Binding{k.Tab, k.Enter, k.Quit}
	case StateResults:
		return []key.Binding{k.Enter, k.Back, k.Quit}
	case StateWeather:
		bindings := []key.Binding{k.Back, k.Quit}
		if k.AddFavorite.Enabled() {
			bindings = append(bindings, k.AddFavorite)
		}
		if k.RemoveFavorite.Enabled() {
			bindings = append(bindings, k.RemoveFavorite)
		}
		return bindings
	default:
		return []key.Binding{k.Quit}
	}
}

// FullHelp returns keybindings for the expanded help view
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

// NewKeyMap creates a new keymap with default bindings
func NewKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "retour"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "changer de focus"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("entrée", "sélectionner"),
		),

		// App control
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quitter"),
		),

		// Favorites management
		AddFavorite: key.NewBinding(
			key.WithKeys("f2"),
			key.WithHelp("f2", "ajouter aux favoris"),
		),
		RemoveFavorite: key.NewBinding(
			key.WithKeys("f3"),
			key.WithHelp("f3", "retirer des favoris"),
		),
	}
}

// UpdateAddRemoveFavoriteBindings updates the enabled state of add/remove favorite keys
func (k *KeyMap) UpdateAddRemoveFavoriteBindings(isFavorite bool) {
	if isFavorite {
		k.AddFavorite.SetEnabled(false)
		k.RemoveFavorite.SetEnabled(true)
	} else {
		k.AddFavorite.SetEnabled(true)
		k.RemoveFavorite.SetEnabled(false)
	}
}

// SetState updates the current application state
func (k *KeyMap) SetState(state ApplicationState) {
	k.State = state
}
