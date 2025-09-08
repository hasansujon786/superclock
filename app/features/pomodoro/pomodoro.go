package pomodoro

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/ui"
	"github.com/hasan/superclock/app/utils"
	"github.com/hasan/superclock/pkg/logger"
)

type PomodoroModel struct {
	timer timer.Model

	modes       []pomoPlayMode
	currentView pomodoroView

	cursor         int
	Running        bool
	Width, Height  int
	CurrentAppView constants.AppView
}

func NewPomodoroModel() PomodoroModel {
	return PomodoroModel{
		currentView: viewDashboard,
		timer:       timer.New(0),
		modes:       getDefaultPlayMode(),
	}
}
func NewPomodoroWithState(ds models.DaemonStateMsg) PomodoroModel {
	logger.Info("NewPomodoroWithState xxxxxxxxxxx")
	return PomodoroModel{
		cursor:      ds.Pomodoro.ModeIdx,
		currentView: utils.If(ds.Pomodoro.Running, viewTimer, viewDashboard),
		timer:       timer.New(ds.Pomodoro.Timeout),
		modes:       getDefaultPlayMode(),
	}
}
func getDefaultPlayMode() []pomoPlayMode {
	return []pomoPlayMode{
		{name: "Work", time: 15 * time.Second},
		{name: "Break", time: 5 * time.Second},
	}
}

func (m PomodoroModel) Reset() tea.Cmd {
	return func() tea.Msg {
		utils.DaemonCmd(constants.CmdReset, nil)()
		return TimerResetMsg{}
	}
}

func (m PomodoroModel) Init() tea.Cmd {
	return m.timer.Init()
}

func (m PomodoroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		var cmd tea.Cmd
		m.currentView = viewDashboard
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case TimerResetMsg:
		m.timer.Timeout = m.modes[m.cursor].time
		return m, nil

	case tea.KeyMsg:
		if constants.AppViewPomodoro != m.CurrentAppView  {
			return m, nil
		}
		
		switch m.currentView {
		case viewTimer: // Timer mappings
			switch msg.String() {
			case tea.KeySpace.String(): // Toggle timer
				return m, tea.Batch(
					m.timer.Toggle(),
					utils.DaemonCmd(constants.CmdToggle, m.timer.Timeout),
				)
			case "r": // reset timer
				if m.timer.Timedout() {
					return m, nil
				}
				return m, tea.Batch(m.timer.Stop(), m.Reset())
			case "s": // stop timer
				m.currentView = viewDashboard
				return m, tea.Batch(m.timer.Stop(), m.Reset())
			}

		case viewEdit: // Timer mappings
			switch msg.String() {
			case tea.KeyEnter.String():
				m.currentView = viewDashboard
				return m, nil
			}

		case viewDashboard: // Dashboard mappings
			switch msg.String() {
			case "i": // View time picker
				m.currentView = viewEdit
			case "s", tea.KeyEnter.String(): // Start timer with currently selected mode
				m.timer = timer.New(m.modes[m.cursor].time)
				m.currentView = viewTimer
				return m, tea.Batch(
					m.timer.Start(),
					utils.DaemonCmd(constants.CmdSetTimer, models.CmdSetTimerPayload{
						ModeIdx: m.cursor,
						Timeout: m.timer.Timeout,
						Play:    true,
					}),
				)
			case tea.KeyLeft.String(), "h":
				if m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyRight.String(), "l":
				if m.cursor < len(m.modes)-1 {
					m.cursor++
				}
			}
		}

		// default:
		// return m, nil
	}

	return m, nil
}

func (m PomodoroModel) View() string {

	var (
		cWidth  = styles.ContainerStyle.GetWidth()
		cHeight = styles.ContainerStyle.GetHeight()

		isRunning = m.timer.Running()
		isPaused  = m.currentView == viewTimer && !isRunning

		header  = styles.HeaderStyle.Render("      ⌛Pomodoro     ")
		content = ""
	)

	switch m.currentView {
	case viewTimer:
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			ui.TimerDigit(utils.FormatDuration(m.timer.Timeout), cWidth, constants.NerdFont),
			"",
			"",
			m.modes[m.cursor].name,
		)
	case viewEdit:
		content = "asdf"
	default: // Dashboard view
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			utils.FormatDurationHumanizeStyled(m.modes[m.cursor].time),
			"",
			"",
			buildModeConfirm(m),
		)
	}

	viewBox := styles.ContainerStyle.Render(
		lipgloss.Place(cWidth, cHeight, lipgloss.Center, lipgloss.Center, content),
	)

	footer := buildFooter(isRunning, isPaused, m.currentView)
	return lipgloss.JoinVertical(lipgloss.Center, header, viewBox, footer)
}

func buildModeConfirm(m PomodoroModel) string {
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 4).
		Background(styles.ThemeColors.Muted)

	var out string
	for i, mode := range m.modes {
		style := buttonStyle
		if i == m.cursor {
			style = style.Background(styles.ThemeColors.Secondary).Foreground(styles.ThemeColors.White)
		}
		out += style.Render(mode.name) + " "
	}

	return out
}

func buildFooter(isRunning, isPaused bool, currentView pomodoroView) string {
	footerMsg := "    "

	if isRunning && !isPaused {
		footerMsg += "space: pause"
	} else {
		footerMsg += "space: play "
	}

	if currentView == viewTimer {
		footerMsg += " • s: stop "
	} else {
		footerMsg += " • s: start"
	}

	footerMsg += " • r: reset\n"

	footerMsg += " "

	return styles.FooterStyle.Render(footerMsg)
}
