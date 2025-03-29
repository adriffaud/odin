package main

import (
	"fmt"
	"os"

	"driffaud.fr/odin/pkg/app"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := app.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
