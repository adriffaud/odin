package ui

import (
	"fmt"

	"driffaud.fr/odin/internal/i18n"
	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/lipgloss"
)

// RenderLoading renders the loading screen
func RenderLoading(spinnerView string, width, height int) string {
	loadingMessage := fmt.Sprintf("%s %s", spinnerView, i18n.T("app.loading", nil))
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(loadingMessage)
}
