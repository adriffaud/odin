package model

import (
	"fmt"

	"driffaud.fr/odin/pkg/service"
	"driffaud.fr/odin/pkg/types"
	"driffaud.fr/odin/pkg/ui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ApplicationState represents the current state of the application
type ApplicationState string

const (
	StateInput     ApplicationState = "input"
	StateSearching ApplicationState = "searching"
	StateResults   ApplicationState = "results"
	StateLoading ApplicationState = "loading"
)

// Model represents the application model
type Model struct {
	width, height int
	state         ApplicationState
	input         textinput.Model
	placesList    list.Model
	selectedPlace types.Place
	err           error
}

// InitialModel returns the initial application model
func InitialModel() Model {
	return Model{
		state:      StateInput,
		input:      ui.InitInput(),
		placesList: ui.InitResultsList(),
		err:        nil,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles state transitions based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.state == StateResults {
				m.state = StateInput
				m.input.Focus()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEnter:
			if m.state == StateInput {
				query := m.input.Value()
				if query == "" {
					return m, nil
				}
				m.state = StateSearching
				return m, service.SearchPlaces(query)
			} else if m.state == StateResults {
				if i, ok := m.placesList.SelectedItem().(types.Place); ok {
					m.selectedPlace = i
					m.state = StateLoading
					return m, service.GetWeather(i.Latitude, i.Longitude)
				}
				return m, tea.Quit
			}
		}

	case service.ErrMsg:
		m.err = msg
		m.state = StateInput
		return m, nil

	case types.SearchResultsMsg:
		m.state = StateResults
		m.placesList.SetItems(msg)
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		h := msg.Height - 6
		if h < 0 {
			h = 10
		}
		m.placesList.SetSize(msg.Width-4, h)
	}

	// Update active component based on state
	if m.state == StateInput || m.state == StateSearching {
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	} else if m.state == StateResults {
		var listCmd tea.Cmd
		m.placesList, listCmd = m.placesList.Update(msg)
		return m, listCmd
	}

	return m, nil
}

// View renders the UI based on the current state
func (m Model) View() string {
	if m.err != nil {
		errorMsg := fmt.Sprintf("Erreur: %s\n\nAppuyer sur une touche pour continuer...", m.err)
		return ui.BorderStyle.
			Width(m.width-2).
			Height(m.height-2).
			Align(lipgloss.Center, lipgloss.Center).
			Render(errorMsg)
	}

	switch m.state {
	case StateInput:
		return ui.InputView(m.input, m.width, m.height)
	case StateSearching:
		return ui.BorderStyle.
			Width(m.width-2).
			Height(m.height-2).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Recherche en cours...")
	case StateResults:
		return ui.ResultsView(m.placesList, m.width, m.height)
	default:
		return ui.BorderStyle.
			Width(m.width-4).
			Height(m.height-4).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Chargement...")
	}
}
