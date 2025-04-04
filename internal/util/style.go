package util

import "github.com/charmbracelet/lipgloss"

var (
	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63"))

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	InputContainerStyle = lipgloss.NewStyle().
				Padding(1, 2)

	WeatherInfoStyle = lipgloss.NewStyle().
				MarginLeft(4)

	WeatherSectionStyle = lipgloss.NewStyle().
				MarginTop(1).
				MarginBottom(1)

	TableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	AstroInfoStyle = lipgloss.NewStyle().
			MarginTop(1).
			Foreground(lipgloss.Color("105")).
			Bold(true)
)
