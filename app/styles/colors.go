package styles

import "github.com/charmbracelet/lipgloss"

type Colors struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Error     lipgloss.Color
	Success   lipgloss.Color
	Muted     lipgloss.Color
	Black     lipgloss.Color
}

var ThemeColors = Colors{
	Primary:   lipgloss.Color("205"), // bright pink
	Secondary: lipgloss.Color("63"),  // Purple
	Accent:    lipgloss.Color("45"),  // teal
	Error:     lipgloss.Color("9"),   // red
	Success:   lipgloss.Color("2"),   // green
	Muted:     lipgloss.Color("240"), // gray
	Black:     lipgloss.Color("#222222"), // gray
}
