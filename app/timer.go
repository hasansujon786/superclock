package app

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/components"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/utils"
)

type TimerClockModel struct {
	timer         timer.Model
	paused        bool
	picker        components.TimeWheelModel
	width, height int
}

func CreateTimerClockModel() TimerClockModel {
	return TimerClockModel{
		timer:  timer.NewWithInterval(0, time.Millisecond),
		picker: components.CreateTimeWheelModel(components.CursorPosSecond),
	}
}

func (m TimerClockModel) Init() tea.Cmd {
	// return m.timer.Init()
	return nil
}

func (m TimerClockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.TimeoutMsg:
		m.paused = false
		m.picker.FocusLast()
		return m, nil

	case tea.KeyMsg:
		switch (constants.ClockState{Running: m.timer.Running(), Paused: m.paused}) {
		// Clock is running or the clock is paused
		case constants.ClockStateRunning, constants.ClockStatePaused:
			switch msg.String() {

			// Reset current timer
			case "r":
				m.timer.Timeout = m.picker.Value.ToDuration()
				return m, nil

				// FIXME: fix layout jumps while pausing
				// Pause and Start timer
			case " ":
				m.paused = !m.paused
				if m.paused {
					return m, m.timer.Stop()
				}
				return m, m.timer.Start()

				// Stop timer & focus input
			case "s", "i", "esc":
				m.paused = false
				m.picker.FocusLast()
				return m, m.timer.Stop()
			}
			
		// CASE: Inside picker
		default:
			var cmd tea.Cmd
			m, cmd = m.pickerKeymaps(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m TimerClockModel) View() string {
	cWidth := styles.ContainerStyle.GetWidth()
	cHeight := styles.ContainerStyle.GetHeight()

	clkState := constants.ClockState{Running: m.timer.Running(), Paused: m.paused}
	isPaused := clkState.IsPaused()
	isRunning := clkState.IsRunning()

	playBtn := components.ButtonStyles.Render("   ")
	escBtn := components.ButtonStyles.Render("   ")
	playPauseBtn := components.ButtonStyles.Render(utils.If(isPaused, "   ", "   "))

	timerContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 5)

	content := ""

	if isRunning || isPaused {
		timeDigit := components.TimerDigit(utils.FormatTimerFromSeconds(m.timer.Timeout), cWidth, components.NerdFont)
		totalTime := utils.FormatTimerFromSeconds(m.picker.Value.ToDuration())

		content = lipgloss.JoinVertical(
			lipgloss.Center,
			timeDigit,
			"",
			totalTime,
			"",
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				playPauseBtn,
				"  ",
				escBtn,
			),
			"",
		)
	} else {
		pickerTime := components.TimerWhell(m.picker.Value, m.picker.Position)
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			"󰀠 Select time ",
			timerContainer.Render(pickerTime),
			"",
			playBtn,
		)
	}

	header := styles.HeaderStyle.Render("      ⌛Timer     ")
	viewBox := styles.ContainerStyle.Render(
		lipgloss.Place(cWidth, cHeight, lipgloss.Center, lipgloss.Center, content),
	)

	footer := buildFooter(isRunning, isPaused)
	return lipgloss.JoinVertical(lipgloss.Center, header, viewBox, footer)
}

func (m TimerClockModel) pickerKeymaps(msg tea.Msg) (TimerClockModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			m.picker.FocusLast()
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
			m.picker.Blur()
		case " ", "s", "enter":
			m.timer.Timeout = m.picker.Value.ToDuration()
			return m, m.timer.Start()
		}
		return m, nil
	}

	return m, nil
}

func buildFooter(isRunning, isPaused bool) string {
	footerMsg := "    "

	if isRunning && !isPaused {
		footerMsg += "space: pause"
	} else {
		footerMsg += "space: play "
	}

	if isRunning {
		footerMsg += " • s: stop "
	} else {
		footerMsg += " • s: start"
	}
	footerMsg += " • r: reset\n"

	footerMsg += "k/j: increase/decrease • h/l: change focus\n"

	return styles.FooterStyle.Render(footerMsg)
}
