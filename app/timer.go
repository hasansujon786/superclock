package app

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/components"
	"github.com/hasan/superclock/app/styles"
)

// --- Timer model ---
type TimerClockModel struct {
	value  int
	timer  timer.Model
	picker components.TimeWheelModel

	width, height int
}

func CreateTimerClockModel() TimerClockModel {
	return TimerClockModel{
		picker: components.CreateTimeWheelModel(),
	}
}

func (m TimerClockModel) Init() tea.Cmd {
	return nil
}

func (m TimerClockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.picker.Position {

		// -----------------------
		// CASE: Nothing selected
		// -----------------------
		case components.CursorPosNone:
			switch msg.String() {
			case "i":
				m.picker.Focus(components.CursorPosSecond) // enter input mode
			}
			return m, nil

		// -----------------------
		// CASE: Inside picker
		// -----------------------
		default:
			switch msg.String() {
			case "h":
				m.picker.PickerMoveCursorLeft()
			case "l":
				m.picker.PickerMoveCursorRight()
			case "k":
				m.picker.IncreaseValue()
			case "j":
				m.picker.DecreaseValue()
			case "r":
				m.picker.Reset()
			case "x":
				m.picker.ResetCurrent()
			case "esc":
				m.picker.Blur() // exit input mode
			case "enter":
				m.picker.Blur() // exit input mode
			}
			return m, nil
		}
	}
	return m, nil
}

func (m TimerClockModel) View() string {
	contentWidth := styles.ContainerStyle.GetWidth() - 2

	time := components.TimerWhell(m.picker.Value, m.picker.Position)
	footer := styles.FooterStyle.Render("↑/↓ to change count • q to quit")

	header := styles.HeaderStyle.Render("      ⌛Timer     ")

	timerContainer := lipgloss.NewStyle().
		Width(contentWidth).
		Height(1).
		MarginTop(1).
		Border(lipgloss.RoundedBorder()).
		Align(lipgloss.Center)

	content := timerContainer.Render(time)

	box := styles.ContainerStyle.Render(content)

	return lipgloss.JoinVertical(lipgloss.Center, header, box, footer)
}
