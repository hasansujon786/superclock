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

type pomodoroView int

const (
	viewDashboard pomodoroView = iota
	viewTimer
	viewEdit
)

type PomodoroModel struct {
	choices     []string
	choiceTimes []time.Duration
	currentView pomodoroView

	timer timer.Model

	cursor        int
	Running       bool
	Width, Height int
}

func NewPomodoroModel() PomodoroModel {
	return PomodoroModel{
		currentView: viewDashboard,
		timer:       timer.New(0),
		choices:     []string{"Work", "Break"},
		choiceTimes: []time.Duration{15 * time.Second, 5 * time.Second},
	}
}
func NewPomodoroWithState(ds models.DaemonStateMsg) PomodoroModel {
	logger.Info("NewPomodoroWithState")
	return PomodoroModel{
		currentView: viewTimer,
		timer:       timer.New(ds.Timeout - ds.Elapsed),
		choices:     []string{"Work", "Break"},
		choiceTimes: []time.Duration{15 * time.Second, 5 * time.Second},
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

	case tea.KeyMsg:
		switch m.currentView {
		case viewDashboard: // Dashboard mappings
			switch msg.String() {
			case "left", "h":
				if m.cursor > 0 {
					m.cursor--
				}
			case "right", "l":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				m.timer.Timeout = m.choiceTimes[m.cursor]
				m.currentView = viewTimer

				return m, tea.Batch(
					m.timer.Start(),
					utils.DaemonCmd(constants.CmdSetTimer, models.CmdSetTimerPayload{
						Timeout: m.timer.Timeout,
						Play:    true,
					}),
				)
			}
		}

		// default:
		// return m, cmd
	}

	return m, nil
}

func (m PomodoroModel) View() string {
	cWidth := styles.ContainerStyle.GetWidth()
	cHeight := styles.ContainerStyle.GetHeight()
	header := styles.HeaderStyle.Render("      ⌛Pomodoro     ")

	// isRunning := m.timer.Running()

	content := ""

	if m.currentView == viewTimer {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			ui.TimerDigit(utils.FormatDuration(m.timer.Timeout), cWidth, constants.NerdFont),
			"",
			"",
		)
	} else {
		content = lipgloss.JoinVertical(
			lipgloss.Center,
			utils.FormatDurationHumanize(m.choiceTimes[m.cursor]),
			"",
			buildModeConfirm(m),
		)
	}

	viewBox := styles.ContainerStyle.Render(
		lipgloss.Place(cWidth, cHeight, lipgloss.Center, lipgloss.Center, content),
	)
	footer := "----"
	return lipgloss.JoinVertical(lipgloss.Center, header, viewBox, footer)
}

func buildModeConfirm(m PomodoroModel) string {
	buttonStyle := lipgloss.NewStyle().
		Padding(0, 4).
		Background(styles.ThemeColors.Muted)

	var out string
	for i, choice := range m.choices {
		style := buttonStyle
		if i == m.cursor {
			style = style.Background(styles.ThemeColors.Secondary).Foreground(styles.ThemeColors.White)
		}
		out += style.Render(choice) + " "
	}

	return out
}
