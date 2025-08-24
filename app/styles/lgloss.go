package styles

import "github.com/charmbracelet/lipgloss"

var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(ThemeColors.Primary).
	Border(lipgloss.RoundedBorder(), false, false, true).
	BorderForeground(ThemeColors.Secondary).
	MarginBottom(1)

var FooterStyle = lipgloss.NewStyle().
	Foreground(ThemeColors.Muted).
	MarginTop(1)

var ContainerStyle = lipgloss.NewStyle().
	Width(40).
	Height(12).
	Border(lipgloss.RoundedBorder()).
	BorderForeground(ThemeColors.Secondary)
