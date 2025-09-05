package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

type App struct {
	view           AppView
	quitting       bool
	timerModel     TimerClockModel
	stopWatchModel StopWatchModel
	pomodoroState  PomodoroModel

	width, height int
}

func NewApp(view AppView, daemonState any) App {
	if data, ok := daemonState.(DaemonStateMsg); ok {
		logger.Info(daemonState)
		logger.Info("NewApp Started...")

		if data.Running {
			return App{
				view:           view,
				timerModel:     NewTimerClockModel(),
				stopWatchModel: NewStopWatchModel(),
				pomodoroState:  NewPomodoroWithState(data),
			}
		}
	}

	return App{
		view:           view,
		timerModel:     NewTimerClockModel(),
		stopWatchModel: NewStopWatchModel(),
		pomodoroState:  NewPomodoroModel(),
	}
}

func (a App) Init() tea.Cmd {
	switch a.view {
	case AppViewTimer:
		return a.timerModel.Init()
	case AppViewStopWatch:
		return a.stopWatchModel.Init()
	case AppViewPomodoro:
		return a.pomodoroState.Init()
	}
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width, a.height = msg.Width, msg.Height

		a.stopWatchModel.height = msg.Height
		a.stopWatchModel.width = msg.Width

		a.timerModel.height = msg.Height
		a.timerModel.width = msg.Width

		a.pomodoroState.height = msg.Height
		a.pomodoroState.width = msg.Width

		utils.NotifyAppMounted()
		return a, nil

	// case DaemonStateMsg:
	// 	logger.Info("DaemonStateMsg")
	// 	return a, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "tab": // Toggle between views
			a.view++
			if a.view >= AppViewNone {
				a.view = AppViewTimer
			}
		case "shift+tab": // Previous view
			a.view--
			if a.view < AppViewTimer {
				a.view = AppViewNone - 1
			}

		case "q", "ctrl+c": // Exit app
			a.quitting = true
			return a, tea.Quit
		}
	}

	// Delegate to current submodel
	switch a.view {
	case AppViewTimer:
		newModel, cmd := a.timerModel.Update(msg)
		a.timerModel = newModel.(TimerClockModel)
		return a, cmd

	case AppViewStopWatch:
		newModel, cmd := a.stopWatchModel.Update(msg)
		a.stopWatchModel = newModel.(StopWatchModel)
		return a, cmd

	case AppViewPomodoro:
		newModel, cmd := a.pomodoroState.Update(msg)
		a.pomodoroState = newModel.(PomodoroModel)
		return a, cmd
	}
	return a, nil
}

func (a App) View() string {
	if a.quitting {
		return ""
	}

	view := ""
	switch a.view {
	case AppViewTimer:
		view = a.timerModel.View()
	case AppViewStopWatch:
		view = a.stopWatchModel.View()
	case AppViewPomodoro:
		view = a.pomodoroState.View()
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center, lipgloss.Center,
		view,
	)
}
