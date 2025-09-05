package timer

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/ui"
	"github.com/hasan/superclock/app/utils"
)

// TimerState matches the daemon TimerState
type TimerState struct {
	Timeout  time.Duration
	Interval time.Duration
	Elapsed  time.Duration
	Running  bool
}

type model struct {
	timer        TimerState
	picker       ui.TimerWheelModel
	err          error
	hasConnError bool
}

func NewModel() model {
	return model{
		picker: ui.NewTimerWheelModel(ui.CursorPosSecond),
	}
}

func (m model) pickerKeymaps(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		case " ", "s", "enter", "esc":
			if m.picker.Value.IsEmpty() {
				return m, nil
			}

			m.picker.Blur()
			return m, tea.Batch(
				sendCmd(constants.CmdSetTimer, m.picker.Value),
				sendCmd(constants.CmdPlay, nil),
			)
		}

	}
	return m, nil
}

// Init connects to the daemon and starts periodic updates
func (m model) Init() tea.Cmd {
	return tickDaemon()
}

// Update handles messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch m.picker.Position {
		// CASE: Inside timer
		case ui.CursorPosNone:
			switch msg.String() {
			case " ":
				// TODO: check time before play
				return m, sendCmd(constants.CmdToggle, nil)
			case "q", "ctrl+c":
				return m, tea.Quit
			case "s", "i":
				m.picker.FocusLast()
				return m, sendCmd(constants.CmdReset, nil)
			case "r":
				return m, sendCmd(constants.CmdReset, nil)
				// case "g":
				// 	return m, sendCmd(constants.CmdGet, nil)
			}

		// CASE: Inside picker
		default:
			var cmd tea.Cmd
			m, cmd = m.pickerKeymaps(msg)
			return m, cmd
		}

	case TimerState:
		m.timer = msg
		// logger.Info("TimerState",)
		return m, tickDaemon() // schedule next update

	case error:
		m.err = msg
		m.hasConnError = true
		return m, tickDaemon()
	}

	return m, nil
}

// View renders the timer
func (m model) View() string {
	if m.hasConnError {
		return fmt.Sprintf("Error connecting to daemon: %v\nPress q to quit", m.err)
	}

	cWidth := styles.ContainerStyle.GetWidth()
	cHeight := styles.ContainerStyle.GetHeight()

	isTimerActive := m.picker.Position == ui.CursorPosNone
	isPaused := !m.timer.Running && isTimerActive

	playBtn := ui.ButtonStyles.Render("   ")
	escBtn := ui.ButtonStyles.Render("   ")
	playPauseBtn := ui.ButtonStyles.Render(utils.If(isPaused, "   ", "   "))

	timerContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 5)

	content := ""

	if isTimerActive {
		timeDigit := ui.TimerDigit(utils.FormatDuration(m.timer.Elapsed), cWidth, constants.NerdFont)
		totalTime := utils.FormatDuration(m.timer.Timeout)

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
		pickerTime := ui.TimerWhell(m.picker.Value, m.picker.Position)
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

	footer := buildFooter(isTimerActive, isPaused)
	return lipgloss.JoinVertical(lipgloss.Center, header, viewBox, footer)
}

// tickDaemon periodically fetches the latest state
func tickDaemon() tea.Cmd {
	return func() tea.Msg {
		return sendCmd(constants.CmdGet, nil)()
	}
}

func buildFooter(isTimerActive, isPaused bool) string {
	footerMsg := "    "

	if isTimerActive && isPaused {
		footerMsg += "space: play "
	} else {
		footerMsg += "space: pause"
	}

	if isTimerActive {
		footerMsg += " • s: stop "
	} else {
		footerMsg += " • s: start"
	}
	footerMsg += " • r: reset\n"

	footerMsg += "k/j: increase/decrease • h/l: change focus\n"

	return styles.FooterStyle.Render(footerMsg)
}
