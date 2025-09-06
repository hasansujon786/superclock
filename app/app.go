package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/features/pomodoro"
	"github.com/hasan/superclock/app/features/stopwatch"
	"github.com/hasan/superclock/app/features/timer"
	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/app/utils"
	"github.com/hasan/superclock/pkg/logger"
)

type AppView int

const (
	AppViewTimer AppView = iota
	AppViewStopWatch
	AppViewPomodoro
	AppViewNone
)

type app struct {
	timer     timer.TimerClockModel
	stopwatch stopwatch.StopWatchModel
	pomodoro  pomodoro.PomodoroModel

	view          AppView
	quitting      bool
	width, height int
}

func NewApp(view AppView, daemonState any) app {
	if data, ok := daemonState.(models.DaemonStateMsg); ok {
		logger.Info(daemonState)
		logger.Info("NewApp Started...")

		if data.Running {
			return app{
				view:      view,
				timer:     timer.NewTimerClockModel(),
				stopwatch: stopwatch.NewStopWatchModel(),
				pomodoro:  pomodoro.NewPomodoroWithState(data),
			}
		}
	}

	return app{
		view:      view,
		timer:     timer.NewTimerClockModel(),
		stopwatch: stopwatch.NewStopWatchModel(),
		pomodoro:  pomodoro.NewPomodoroModel(),
	}
}

func (a app) Init() tea.Cmd {
	switch a.view {
	case AppViewTimer:
		return a.timer.Init()
	case AppViewStopWatch:
		return a.stopwatch.Init()
	case AppViewPomodoro:
		return a.pomodoro.Init()
	}
	return nil
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height
		utils.NotifyAppMounted()
		return a.updateSubModels(msg) // Pass resize msg to all sub models too

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			a.view++
			if a.view >= AppViewNone {
				a.view = AppViewTimer
			}
		case "shift+tab":
			a.view--
			if a.view < AppViewTimer {
				a.view = AppViewNone - 1
			}
		case "q", "ctrl+c":
			a.quitting = true
			return a, tea.Quit
		}
	}

	return a.updateSubModels(msg)
}

func (a app) updateSubModels(msg tea.Msg) (app, tea.Cmd) {
	var cmds []tea.Cmd

	// Timer
	newTimer, cmd := a.timer.Update(msg)
	a.timer = newTimer.(timer.TimerClockModel)
	cmds = append(cmds, cmd)

	// Stopwatch
	newStopwatch, cmd := a.stopwatch.Update(msg)
	a.stopwatch = newStopwatch.(stopwatch.StopWatchModel)
	cmds = append(cmds, cmd)

	// Pomodoro
	newPomodoro, cmd := a.pomodoro.Update(msg)
	a.pomodoro = newPomodoro.(pomodoro.PomodoroModel)
	cmds = append(cmds, cmd)

	// Only render active view, but keep others updated
	return a, tea.Batch(cmds...)
}

func (a app) View() string {
	if a.quitting {
		return ""
	}

	view := ""
	switch a.view {
	case AppViewTimer:
		view = a.timer.View()
	case AppViewStopWatch:
		view = a.stopwatch.View()
	case AppViewPomodoro:
		view = a.pomodoro.View()
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center, lipgloss.Center,
		view,
	)
}
