package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/utils"
)

type AppView int

const (
	timerClockView AppView = iota
	stopWatchView
)

type App struct {
	view          AppView
	quitting       bool
	timerModel     TimerClockModel
	stopWatchModel StopWatchModel

	width, height int
}

func NewApp() App {
	return App{
		view:          timerClockView, // start on timer screen
		timerModel:     CreateTimerClockModel(),
		stopWatchModel: DefaultStopWatchModel(),
	}
}

func (a App) Init() tea.Cmd {
	// Delegate to the current submodel
	switch a.view {
	case timerClockView:
		return a.timerModel.Init()
	case stopWatchView:
		return a.stopWatchModel.Init()
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
		utils.NotifyAppMounted()
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "tab": // Toggle between views
			if a.view == timerClockView {
				a.view = stopWatchView
			} else {
				a.view = timerClockView
			}

		case "q", "ctrl+c": // Exit app
			a.quitting = true
			return a, tea.Quit
		}
	}

	// Delegate to current submodel
	switch a.view {
	case timerClockView:
		newModel, cmd := a.timerModel.Update(msg)
		a.timerModel = newModel.(TimerClockModel)
		return a, cmd

	case stopWatchView:
		newModel, cmd := a.stopWatchModel.Update(msg)
		a.stopWatchModel = newModel.(StopWatchModel)
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
	case timerClockView:
		view = a.timerModel.View()
	case stopWatchView:
		view = a.stopWatchModel.View()
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center, lipgloss.Center,
		view,
	)
}
