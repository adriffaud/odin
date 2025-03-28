package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// InitInput initializes the text input component
func InitInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "Entrer un nom de lieu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40
	return ti
}

// InputView renders the input screen
func InputView(input textinput.Model, width, height int) string {
	title := TitleStyle.Render("Entrer un nom de lieu")
	inputField := input.View()
	hint := "(entrer pour rechercher, esc pour quitter)"

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		inputField,
		"",
		hint,
	)

	styled := InputContainerStyle.Render(content)
	return BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(styled)
}
