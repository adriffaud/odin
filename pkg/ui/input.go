package ui

import "github.com/charmbracelet/bubbles/textinput"

// InitInput initializes the text input component
func InitInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Entrer un nom de lieu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return ti
}

// InputView renders the input screen
func InputView(input textinput.Model) string {
	return "Entrer un nom de lieu :\n\n" + input.View() + "\n\n(enter to search, esc to quit)"
}
