package main

import (
	"fmt"
	"os"

	"driffaud.fr/odin/internal/app"
	"driffaud.fr/odin/internal/i18n"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if err := i18n.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "error initializing i18n: %v\n", err)
		os.Exit(1)
	}

	model := app.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running program: %v\n", err)
		os.Exit(1)
	}
}
