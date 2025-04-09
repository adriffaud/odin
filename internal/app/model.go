package app

import (
	"time"

	"driffaud.fr/odin/internal/app/ui"
	"driffaud.fr/odin/internal/domain"
	"driffaud.fr/odin/internal/platform/api/openmeteo"
	"driffaud.fr/odin/internal/platform/api/photon"
	"driffaud.fr/odin/internal/platform/storage"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	placeModel    ui.PlaceModel
	weatherModel  ui.WeatherModel
	placesList    list.Model
	weatherData   domain.WeatherData
	selectedPlace domain.Place
	spinner       spinner.Model
	favorites     *storage.FavoritesStore
	err           error
	keyMap        KeyMap
	help          help.Model
}

type placesResultMsg struct {
	places []domain.Place
	err    error
}

type weatherResultMsg struct {
	data domain.WeatherData
	err  error
}

// InitialModel returns the initial application model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Moon
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	favStore, err := storage.NewFavoritesStore()
	if err != nil {
		favStore = &storage.FavoritesStore{}
	}

	placeModel := ui.NewPlaceModel(favStore)
	helpModel := help.New()
	helpModel.ShowAll = false

	return Model{
		state:      StatePlace,
		placeModel: placeModel,
		placesList: ui.InitResultsList(),
		spinner:    s,
		favorites:  favStore,
		err:        nil,
		keyMap:     NewKeyMap(),
		help:       helpModel,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Sequence(
		m.placeModel.Init(),
		m.weatherModel.Init(),
		tea.SetWindowTitle("Odin"),
	)
}

// Update handles state transitions based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case placesResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StatePlace
			return m, nil
		}
		m.state = StateResults
		items := make([]list.Item, len(msg.places))
		for i, p := range msg.places {
			items[i] = p
		}
		m.placesList.SetItems(items)
		return m, nil
	case weatherResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = StatePlace
			return m, nil
		}
		return m.handleWeatherResultMsg(msg.data)
	case tea.WindowSizeMsg:
		return m.handleWindowSizeMsg(msg)
	}

	return m.updateActiveComponent(msg)
}

// View renders the UI based on the current state
func (m Model) View() string {
	if m.err != nil {
		return ui.RenderError(m.err, m.width, m.height)
	}

	m.keyMap.SetState(m.state)
	m.help.ShowAll = m.state == StateWeather
	helpView := m.help.View(m.keyMap)

	switch m.state {
	case StatePlace:
		return m.placeModel.View(helpView)
	case StateLoading:
		return ui.RenderLoading(m.spinner.View(), m.width, m.height)
	case StateResults:
		return ui.RenderResults(m.placesList, helpView, m.width, m.height)
	case StateWeather:
		return m.weatherModel.View(helpView)
	default:
		return ui.RenderLoading(m.spinner.View(), m.width, m.height)
	}
}

func (m Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keyMap.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keyMap.Back):
		if m.state == StateResults || m.state == StateWeather {
			m.state = StatePlace
		}
		return m, nil
	case key.Matches(msg, m.keyMap.Enter):
		return m.handleEnterKey()
	case key.Matches(msg, m.keyMap.AddFavorite):
		return m.handleAddFavorite()
	case key.Matches(msg, m.keyMap.RemoveFavorite):
		return m.handleRemoveFavorite()
	}

	return m.updateActiveComponent(msg)
}

func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	switch m.state {
	case StatePlace:
		focus := m.placeModel.GetFocusIndex()
		switch focus {
		case 0:
			query := m.placeModel.GetQuery()
			if query == "" {
				return m, nil
			}
			m.state = StateLoading
			cmd := func() tea.Msg {
				places, err := photon.SearchPlaces(query)
				time.Sleep(500 * time.Millisecond)
				return placesResultMsg{places: places, err: err}
			}
			return m, tea.Batch(cmd, m.spinner.Tick)
		case 1:
			if place, ok := m.placeModel.GetSelectedFavorite(); ok {
				m.selectedPlace = place
				m.state = StateLoading
				cmd := func() tea.Msg {
					weather, err := openmeteo.GetWeather(place.Latitude, place.Longitude)
					time.Sleep(500 * time.Millisecond)
					return weatherResultMsg{data: weather, err: err}
				}
				return m, tea.Batch(cmd, m.spinner.Tick)
			}
		}
	case StateResults:
		if i, ok := m.placesList.SelectedItem().(domain.Place); ok {
			m.selectedPlace = i
			m.state = StateLoading
			cmd := func() tea.Msg {
				weather, err := openmeteo.GetWeather(i.Latitude, i.Longitude)
				return weatherResultMsg{data: weather, err: err}
			}
			return m, tea.Batch(cmd, m.spinner.Tick)
		}
	}
	return m, nil
}

func (m Model) handleAddFavorite() (tea.Model, tea.Cmd) {
	if m.state == StateWeather && !m.favorites.IsFavorite(m.selectedPlace) {
		if err := m.favorites.AddFavorite(m.selectedPlace); err != nil {
			m.err = err
			return m, nil
		}
		m.keyMap.UpdateAddRemoveFavoriteBindings(true)
		m.placeModel.UpdateFavorites()
		return m, nil
	}
	return m, nil
}

func (m Model) handleRemoveFavorite() (tea.Model, tea.Cmd) {
	if m.state == StateWeather && m.favorites.IsFavorite(m.selectedPlace) {
		if err := m.favorites.RemoveFavorite(m.selectedPlace); err != nil {
			m.err = err
			return m, nil
		}
		m.keyMap.UpdateAddRemoveFavoriteBindings(false)
		m.placeModel.UpdateFavorites()
		return m, nil
	}
	return m, nil
}

func (m Model) handleWeatherResultMsg(data domain.WeatherData) (tea.Model, tea.Cmd) {
	m.weatherData = data
	m.state = StateWeather
	m.weatherModel = ui.NewWeatherModel(
		data,
		m.selectedPlace,
		m.favorites,
		m.width,
		m.height,
	)
	isFavorite := m.favorites.IsFavorite(m.selectedPlace)
	m.keyMap.UpdateAddRemoveFavoriteBindings(isFavorite)
	return m, m.weatherModel.Init()
}

func (m Model) handleWindowSizeMsg(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	h := msg.Height - 6
	if h < 0 {
		h = 10
	}
	m.placesList.SetSize(msg.Width-4, h)

	return m.updateActiveComponent(msg)
}

func (m Model) updateActiveComponent(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StatePlace:
		var placeCmd tea.Cmd
		m.placeModel, placeCmd = m.placeModel.Update(msg)
		return m, placeCmd
	case StateLoading:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case StateResults:
		var listCmd tea.Cmd
		m.placesList, listCmd = m.placesList.Update(msg)
		return m, listCmd
	case StateWeather:
		var weatherCmd tea.Cmd
		isFavorite := m.favorites.IsFavorite(m.selectedPlace)
		m.keyMap.UpdateAddRemoveFavoriteBindings(isFavorite)
		m.weatherModel, weatherCmd = m.weatherModel.Update(msg)
		return m, weatherCmd
	}
	return m, nil
}
