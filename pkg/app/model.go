package app

import (
	"driffaud.fr/odin/pkg/weather"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ApplicationState represents the current state of the application
type ApplicationState string

const (
	StatePlace   ApplicationState = "place"
	StateResults ApplicationState = "results"
	StateWeather ApplicationState = "weather"
	StateLoading ApplicationState = "loading"
)

// Model represents the application model
type Model struct {
	width, height int
	state         ApplicationState
	placeModel    weather.PlaceModel
	weatherModel  weather.WeatherModel
	placesList    list.Model
	weatherData   weather.WeatherData
	selectedPlace weather.Place
	spinner       spinner.Model
	favorites     *weather.FavoritesStore
	err           error
}

// InitialModel returns the initial application model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	favStore, err := weather.NewFavoritesStore()
	if err != nil {
		favStore = &weather.FavoritesStore{}
	}

	placeModel := weather.NewPlaceModel(favStore)

	return Model{
		state:      StatePlace,
		placeModel: placeModel,
		placesList: InitResultsList(),
		spinner:    s,
		favorites:  favStore,
		err:        nil,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.placeModel.Init(),
		m.weatherModel.Init(),
		tea.SetWindowTitle("Odin"),
	)
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
			if m.state == StateResults || m.state == StateWeather {
				m.state = StatePlace
			}
			return m, nil

		case tea.KeyEnter:
			if m.state == StatePlace {
				focus := m.placeModel.GetFocusIndex()
				if focus == 0 {
					query := m.placeModel.GetQuery()
					if query == "" {
						return m, nil
					}
					m.state = StateLoading
					return m, tea.Batch(
						weather.SearchPlaces(query),
						m.spinner.Tick,
					)
				} else if focus == 1 {
					if place, ok := m.placeModel.GetSelectedFavorite(); ok {
						m.selectedPlace = place
						m.state = StateLoading
						return m, tea.Batch(
							weather.GetWeather(place.Latitude, place.Longitude),
							m.spinner.Tick,
						)
					}
				}
			} else if m.state == StateResults {
				if i, ok := m.placesList.SelectedItem().(weather.Place); ok {
					m.selectedPlace = i
					m.state = StateLoading
					return m, tea.Batch(
						weather.GetWeather(i.Latitude, i.Longitude),
						m.spinner.Tick,
					)
				}
				return m, tea.Quit
			}
		}

	case weather.ErrMsg:
		m.err = msg
		m.state = StatePlace
		return m, nil

	case weather.SearchResultsMsg:
		m.state = StateResults
		m.placesList.SetItems(msg)
		return m, nil

	case weather.WeatherResultMsg:
		m.weatherData = msg.Data
		m.state = StateWeather
		m.weatherModel = weather.NewWeatherModel(
			msg.Data,
			m.selectedPlace,
			m.favorites,
			m.width,
			m.height,
		)
		return m, m.weatherModel.Init()

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
	switch m.state {
	case StatePlace:
		var placeCmd tea.Cmd
		m.placeModel, placeCmd = m.placeModel.Update(msg)
		return m, placeCmd
	case StateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	case StateResults:
		var listCmd tea.Cmd
		m.placesList, listCmd = m.placesList.Update(msg)
		return m, listCmd
	case StateWeather:
		var weatherCmd tea.Cmd
		m.weatherModel, weatherCmd = m.weatherModel.Update(msg)
		return m, weatherCmd
	}

	return m, nil
}

// View renders the UI based on the current state
func (m Model) View() string {
	if m.err != nil {
		return RenderError(m.err, m.width, m.height)
	}

	switch m.state {
	case StatePlace:
		return m.placeModel.View()
	case StateLoading:
		return RenderLoading(m.spinner.View(), m.width, m.height)
	case StateResults:
		return RenderResults(m.placesList, m.width, m.height)
	case StateWeather:
		return m.weatherModel.View()
	default:
		return RenderLoading(m.spinner.View(), m.width, m.height)
	}
}
