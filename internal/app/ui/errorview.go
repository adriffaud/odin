package ui

import (
	"fmt"

	"driffaud.fr/odin/internal/i18n"
	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/lipgloss"
)

// RenderError renders the error screen
func RenderError(err error, width, height int) string {
	errorMsg := fmt.Sprint(i18n.T("app.error", map[string]any{"Error": err}))
	return util.BorderStyle.
		Width(width-2).
		Height(height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(errorMsg)
}
