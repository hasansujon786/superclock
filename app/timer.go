package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Timer model ---
type TimerModel struct {
	value int
	count int

	width, height int
}

func (m TimerModel) Init() tea.Cmd {
	return nil
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "k": // Exit app
			m.count++
			return m, nil

		case "j": // Exit app
			m.count--
			return m, nil
		}

	}

	return m, nil
}

func (m TimerModel) View() string {
	boxStyle := lipgloss.NewStyle().
		Width(40).
		Height(12).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)

	contentWidth := boxStyle.GetWidth() - 2

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Underline(true).
		// Margin(1, 0).
		Width(contentWidth).
		AlignHorizontal(lipgloss.Center)

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	title := titleStyle.Render("hasan")
	footer := footerStyle.Render("↑/↓ to change count • q to quit")

	content := lipgloss.JoinVertical(lipgloss.Left, title)

	box := boxStyle.Render(content)

	return lipgloss.JoinVertical(lipgloss.Center, box, footer)
}
