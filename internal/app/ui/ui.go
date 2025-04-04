package ui

import (
	"fmt"

	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// InitResults initializes the results list component
func InitResultsList() list.Model {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Résultats"
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

// RenderLoading renders the loading screen
func RenderLoading(spinnerView string, width, height int) string {
	loadingMessage := fmt.Sprintf("%s Chargement...", spinnerView)
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(loadingMessage)
}

// RenderError renders the error screen
func RenderError(err error, width, height int) string {
	errorMsg := fmt.Sprintf("Erreur: %s\n\nAppuyer sur une touche pour continuer...", err)
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(errorMsg)
}
