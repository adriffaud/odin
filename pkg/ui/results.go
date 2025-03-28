package ui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// InitResults initializes the results list component
func InitResultsList() list.Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Résultats"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = TitleStyle
	l.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)
	return l
}

// ResultsView renders the results screen
func ResultsView(resultsList list.Model, width, height int) string {
	hint := "(entrer pour sélectionner, esc pour quitter)"
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		resultsList.View(),
		"",
		hint,
	)

	return BorderStyle.
		Width(width - 2).
		Height(height - 2).
		Render(content)
}
