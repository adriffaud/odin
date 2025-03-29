package model

import (
	"fmt"

	"driffaud.fr/odin/pkg/service"
	"driffaud.fr/odin/pkg/types"
	"driffaud.fr/odin/pkg/ui"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ApplicationState represents the current state of the application
type ApplicationState string

const (
	StateInput   ApplicationState = "input"
	StateResults ApplicationState = "results"
	StateWeather ApplicationState = "weather"
	StateLoading ApplicationState = "loading"
)

// Model represents the application model
type Model struct {
	width, height int
	state         ApplicationState
	input         textinput.Model
	placesList    list.Model
	weatherData   types.WeatherData
	selectedPlace types.Place
	spinner       spinner.Model
	err           error
}

// InitialModel returns the initial application model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		state:      StateInput,
		input:      ui.InitInput(),
		placesList: ui.InitResultsList(),
		spinner:    s,
		err:        nil,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.SetWindowTitle("Odin"))
}

// Update handles state transitions based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			if m.state == StateResults {
				m.state = StateInput
				m.input.Focus()
				return m, nil
			} else if m.state == StateWeather {
				m.state = StateResults
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyEnter:
			if m.state == StateInput {
				query := m.input.Value()
				if query == "" {
					return m, nil
				}
				m.state = StateLoading
				return m, tea.Batch(
					service.SearchPlaces(query),
					m.spinner.Tick,
				)
			} else if m.state == StateResults {
				if i, ok := m.placesList.SelectedItem().(types.Place); ok {
					m.selectedPlace = i
					m.state = StateLoading
					return m, tea.Batch(
						service.GetWeather(i.Latitude, i.Longitude),
						m.spinner.Tick,
					)
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

	case types.WeatherResultMsg:
		m.weatherData = msg.Data
		m.state = StateWeather
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
	if m.state == StateInput {
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	} else if m.state == StateLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
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

	loadingMessage := fmt.Sprintf("%s Chargement...", m.spinner.View())
	loadingScreen := ui.BorderStyle.
		Width(m.width-2).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(loadingMessage)

	switch m.state {
	case StateInput:
		return ui.InputView(m.input, m.width, m.height)
	case StateLoading:
		return loadingScreen
	case StateResults:
		return ui.ResultsView(m.placesList, m.width, m.height)
	case StateWeather:
		placeName := m.selectedPlace.Name + " (" + m.selectedPlace.Address + ")"
		weatherContent := ui.WeatherView(m.weatherData, placeName, m.width, m.height)
		return ui.BorderStyle.
			Width(m.width - 2).
			Height(m.height - 2).
			Render(weatherContent)
	default:
		return loadingScreen
	}
}
