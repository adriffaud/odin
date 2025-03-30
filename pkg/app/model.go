package app

import (
	"driffaud.fr/odin/pkg/weather"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
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
	input         textinput.Model
	placesList    list.Model
	favoritesList list.Model
	weatherData   weather.WeatherData
	selectedPlace weather.Place
	spinner       spinner.Model
	favorites     *FavoritesStore
	focusIndex    int // 0 for input, 1 for favorites list
	err           error
}

// InitialModel returns the initial application model
func InitialModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	favStore, err := NewFavoritesStore()
	if err != nil {
		favStore = &FavoritesStore{}
	}

	var favoriteItems []list.Item
	for _, fav := range favStore.Favorites {
		favoriteItems = append(favoriteItems, fav)
	}

	favoritesList := InitFavoritesList(favoriteItems)

	return Model{
		state:         StatePlace,
		input:         InitInput(),
		placesList:    InitResultsList(),
		favoritesList: favoritesList,
		spinner:       s,
		favorites:     favStore,
		focusIndex:    0,
		err:           nil,
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
			if m.state == StateResults || m.state == StateWeather {
				m.state = StatePlace
				m.input.Focus()
				m.focusIndex = 0
				return m, nil
			} else if m.state == StatePlace && m.focusIndex == 1 {
				m.focusIndex = 0
				m.input.Focus()
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyTab:
			if m.state == StatePlace {
				if m.focusIndex == 0 && len(m.favorites.Favorites) > 0 {
					m.focusIndex = 1
					m.input.Blur()
					return m, nil
				} else if m.focusIndex == 1 {
					m.focusIndex = 0
					m.input.Focus()
					return m, nil
				}
			}

		case tea.KeyEnter:
			if m.state == StatePlace {
				if m.focusIndex == 0 {
					query := m.input.Value()
					if query == "" {
						return m, nil
					}
					m.state = StateLoading
					return m, tea.Batch(
						weather.SearchPlaces(query),
						m.spinner.Tick,
					)
				} else if m.focusIndex == 1 {
					if i, ok := m.favoritesList.SelectedItem().(weather.Place); ok {
						m.selectedPlace = i
						m.state = StateLoading
						return m, tea.Batch(
							weather.GetWeather(i.Latitude, i.Longitude),
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

		case tea.KeyF2:
			if m.state == StateWeather {
				m.favorites.AddFavorite(m.selectedPlace)

				// Update favorites list
				var favoriteItems []list.Item
				for _, fav := range m.favorites.Favorites {
					favoriteItems = append(favoriteItems, fav)
				}
				m.favoritesList.SetItems(favoriteItems)

				return m, nil
			}

		case tea.KeyF3:
			if m.state == StateWeather && m.favorites.IsFavorite(m.selectedPlace) {
				m.favorites.RemoveFavorite(m.selectedPlace)

				// Update favorites list
				var favoriteItems []list.Item
				for _, fav := range m.favorites.Favorites {
					favoriteItems = append(favoriteItems, fav)
				}
				m.favoritesList.SetItems(favoriteItems)

				return m, nil
			}
		}

	case weather.ErrMsg:
		m.err = msg
		m.state = StatePlace
		m.focusIndex = 0
		m.input.Focus()
		return m, nil

	case weather.SearchResultsMsg:
		m.state = StateResults
		m.placesList.SetItems(msg)
		return m, nil

	case weather.WeatherResultMsg:
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

		// Set size for favorites list to be about 1/3 of the screen height
		favHeight := max((msg.Height/3)-4, 3)
		m.favoritesList.SetSize(msg.Width-4, favHeight)
	}

	// Update active component based on state
	if m.state == StatePlace {
		if m.focusIndex == 0 {
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		} else {
			var listCmd tea.Cmd
			m.favoritesList, listCmd = m.favoritesList.Update(msg)
			return m, listCmd
		}
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
		return RenderError(m.err, m.width, m.height)
	}

	switch m.state {
	case StatePlace:
		return RenderPlaces(m.input, m.favoritesList, m.focusIndex, m.width, m.height)
	case StateLoading:
		return RenderLoading(m.spinner.View(), m.width, m.height)
	case StateResults:
		return RenderResults(m.placesList, m.width, m.height)
	case StateWeather:
		placeName := m.selectedPlace.Name + " (" + m.selectedPlace.Address + ")"
		isFavorite := m.favorites.IsFavorite(m.selectedPlace)
		return RenderWeather(m.weatherData, placeName, isFavorite, m.width, m.height)
	default:
		return RenderLoading(m.spinner.View(), m.width, m.height)
	}
}
