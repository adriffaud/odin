package ui

import (
	"driffaud.fr/odin/internal/domain"
	"driffaud.fr/odin/internal/i18n"
	"driffaud.fr/odin/internal/platform/storage"
	"driffaud.fr/odin/internal/util"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PlaceModel manages the place search and favorites UI
type PlaceModel struct {
	width, height int
	input         textinput.Model
	favoritesList list.Model
	focusIndex    int // 0 for input, 1 for favorites list
	favorites     *storage.FavoritesStore
}

// NewPlaceModel initializes a new place search model
func NewPlaceModel(favorites *storage.FavoritesStore) PlaceModel {
	ti := textinput.New()
	ti.Placeholder = i18n.T("app.search_place", nil)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	var favoriteItems []list.Item
	for _, fav := range favorites.Favorites {
		favoriteItems = append(favoriteItems, fav)
	}

	favoritesList := list.New(favoriteItems, list.NewDefaultDelegate(), 0, 0)
	favoritesList.Title = i18n.T("app.favorites_title", nil)
	favoritesList.SetShowStatusBar(false)
	favoritesList.SetFilteringEnabled(false)
	favoritesList.SetShowHelp(false)
	favoritesList.Styles.Title = util.SubtitleStyle

	return PlaceModel{
		input:         ti,
		favoritesList: favoritesList,
		favorites:     favorites,
		focusIndex:    0,
	}
}

// Init initializes the model
func (m PlaceModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the place model
func (m PlaceModel) Update(msg tea.Msg) (PlaceModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Set size for favorites list to be about 1/3 of the screen height
		favHeight := max((msg.Height/2)-4, 3)
		m.favoritesList.SetSize(msg.Width-4, favHeight)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
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
	}

	// Handle updates for the focused component
	if m.focusIndex == 0 {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		cmd = inputCmd
	} else {
		var listCmd tea.Cmd
		m.favoritesList, listCmd = m.favoritesList.Update(msg)
		cmd = listCmd
	}

	return m, cmd
}

// View renders the place model UI
func (m PlaceModel) View(helpView string) string {
	title := util.TitleStyle.Render(i18n.T("app.title", nil))

	inputTitle := i18n.T("app.search_title", nil)
	if m.focusIndex == 0 {
		inputTitle = "> " + inputTitle + " <"
	}
	inputTitleStyled := util.SubtitleStyle.Render(inputTitle)
	inputField := lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Render(m.input.View())

	var favoritesSection string

	if len(m.favoritesList.Items()) > 0 {
		if m.focusIndex == 1 {
			m.favoritesList.Title = "> " + m.favoritesList.Title + " <"
		}

		favoritesSection = m.favoritesList.View()
	} else {
		favoritesSection = lipgloss.NewStyle().
			Faint(true).
			Render(i18n.T("app.no_favorites", nil))
	}

	inputSection := lipgloss.JoinVertical(lipgloss.Left,
		inputTitleStyled,
		inputField)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, inputSection, favoritesSection),
		"",
		helpView,
	)

	return util.BorderStyle.
		Width(m.width-2).
		Height(m.height-2).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

// GetQuery returns the current input value
func (m PlaceModel) GetQuery() string {
	return m.input.Value()
}

// GetSelectedFavorite returns the currently selected favorite place if any
func (m PlaceModel) GetSelectedFavorite() (domain.Place, bool) {
	if m.focusIndex == 1 {
		if i, ok := m.favoritesList.SelectedItem().(domain.Place); ok {
			return i, true
		}
	}
	return domain.Place{}, false
}

// GetFocusIndex returns the current focus index
func (m PlaceModel) GetFocusIndex() int {
	return m.focusIndex
}

// UpdateFavorites updates the favorites list with the current favorites
func (m *PlaceModel) UpdateFavorites() {
	var favoriteItems []list.Item
	for _, fav := range m.favorites.Favorites {
		favoriteItems = append(favoriteItems, fav)
	}
	m.favoritesList.SetItems(favoriteItems)
}
