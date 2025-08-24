package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/clock-tui/app/styles"
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
	contentWidth := styles.ContainerStyle.GetWidth() - 2

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.ThemeColors.Primary).
		Underline(true).
		// Margin(1, 0).
		Width(contentWidth).
		AlignHorizontal(lipgloss.Center)

	time := titleStyle.Render(fmt.Sprintf("%d", m.count))
	footer := styles.FooterStyle.Render("↑/↓ to change count • q to quit")

	header := styles.HeaderStyle.Render("      ⌛Timer     ")
	content := lipgloss.JoinVertical(lipgloss.Left, time)

	box := styles.ContainerStyle.Render(content)

	return lipgloss.JoinVertical(lipgloss.Center, header, box, footer)
}
