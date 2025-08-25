package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/utils"
)

type state int

const (
	timer state = iota
	stopWatch
)

type App struct {
	state          state
	quitting       bool
	timerModel     TimerModel
	stopWatchModel StopWatchModel

	width, height int
}

func NewApp() App {
	return App{
		state:          stopWatch, // start on timer screen
		timerModel:     TimerModel{value: 0},
		stopWatchModel: DefaultStopWatchModel(),
	}
}

func (a App) Init() tea.Cmd {
	// Delegate to the current submodel
	switch a.state {
	case timer:
		return a.timerModel.Init()
	case stopWatch:
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
			if a.state == timer {
				a.state = stopWatch
			} else {
				a.state = timer
			}

		case "q", "ctrl+c": // Exit app
			a.quitting = true
			return a, tea.Quit
		}
	}

	// Delegate to current submodel
	switch a.state {
	case timer:
		newModel, cmd := a.timerModel.Update(msg)
		a.timerModel = newModel.(TimerModel)
		return a, cmd

	case stopWatch:
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
	switch a.state {
	case timer:
		view = a.timerModel.View()
	case stopWatch:
		view = a.stopWatchModel.View()
	}

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center, lipgloss.Center,
		view,
	)
}
