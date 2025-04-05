package ui

import (
	"driffaud.fr/odin/internal/i18n"
	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// InitResults initializes the results list component
func InitResultsList() list.Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = i18n.T("app.results_title", nil)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = util.TitleStyle
	l.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)
	return l
}

// RenderResults renders the results screen
func RenderResults(resultsList list.Model, helpView string, width, height int) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		resultsList.View(),
		"",
		helpView,
	)

	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}
