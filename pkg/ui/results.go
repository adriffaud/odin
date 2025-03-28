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
	l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)
	return l
}

// ResultsView renders the results screen
func ResultsView(resultsList list.Model) string {
	return resultsList.View() + "\n\n(enter to select, esc to quit)"
}
