package app

import (
	"fmt"

	"driffaud.fr/odin/pkg/util"
	"github.com/charmbracelet/bubbles/list"
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

// InitFavoritesList initializes the favorites list component
func InitFavoritesList(items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Lieux favoris"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = util.TitleStyle
	l.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)
	return l
}

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
func RenderResults(resultsList list.Model, width, height int) string {
	hint := "(entrer pour sélectionner, esc pour quitter)"
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		resultsList.View(),
		"",
		hint,
	)

	return util.BorderStyle.
		Width(width - 2).
		Height(height - 2).
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
