package weather

import (
	"driffaud.fr/odin/pkg/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Tab   key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Tab, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Tab, k.Enter, k.Quit},
	}
}

var keys = keyMap{
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "changer de focus"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("entrée", "sélectionner"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quitter"),
	),
}

// PlaceModel manages the place search and favorites UI
type PlaceModel struct {
	width, height int
	input         textinput.Model
	favoritesList list.Model
	focusIndex    int // 0 for input, 1 for favorites list
	favorites     *FavoritesStore
	help          help.Model
	keys          keyMap
}

// NewPlaceModel initializes a new place search model
func NewPlaceModel(favorites *FavoritesStore) PlaceModel {
	// Initialize text input
	ti := textinput.New()
	ti.Placeholder = "Entrer un nom de lieu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	// Initialize favorites list
	var favoriteItems []list.Item
	for _, fav := range favorites.Favorites {
		favoriteItems = append(favoriteItems, fav)
	}

	favoritesList := list.New(favoriteItems, list.NewDefaultDelegate(), 0, 0)
	favoritesList.Title = "Lieux favoris"
	favoritesList.SetShowStatusBar(false)
	favoritesList.SetFilteringEnabled(false)
	favoritesList.SetShowHelp(false)
	favoritesList.Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true)
	favoritesList.Styles.HelpStyle = lipgloss.NewStyle().MarginLeft(2)

	help := help.New()
	help.ShowAll = false

	return PlaceModel{
		input:         ti,
		favoritesList: favoritesList,
		favorites:     favorites,
		focusIndex:    0,
		help:          help,
		keys:          keys,
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
func (m PlaceModel) View() string {
	title := util.TitleStyle.Render("Météo astronomique")

	inputTitle := "Rechercher un lieu"
	if m.focusIndex == 0 {
		inputTitle = "> " + inputTitle + " <"
	}
	inputTitleStyled := lipgloss.NewStyle().Bold(true).Render(inputTitle)
	inputField := lipgloss.NewStyle().
		PaddingTop(1).
		PaddingBottom(1).
		Render(m.input.View())

	var favoritesSection string

	if len(m.favoritesList.Items()) > 0 {
		favoritesTitle := "Favoris"
		if m.focusIndex == 1 {
			favoritesTitle = "> " + favoritesTitle + " <"
		}

		favoritesSection = lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Render(favoritesTitle),
			m.favoritesList.View(),
		)
	} else {
		favoritesSection = lipgloss.NewStyle().
			Faint(true).
			Render("Aucun lieu favori - Appuyez sur F2 pour en ajouter")
	}

	helpView := m.help.View(m.keys)

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
func (m PlaceModel) GetSelectedFavorite() (Place, bool) {
	if m.focusIndex == 1 {
		if i, ok := m.favoritesList.SelectedItem().(Place); ok {
			return i, true
		}
	}
	return Place{}, false
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
