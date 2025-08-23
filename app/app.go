package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	tabColors := [2]string{"240", "240"}
	view := ""
	switch a.state {
	case timer:
		tabColors[0] = "63"
		view = a.timerModel.View()
	case stopWatch:
		tabColors[1] = "63"
		view = a.stopWatchModel.View()
	}

	header := ""

	color1 := lipgloss.Color(tabColors[0])
	tabTimer := lipgloss.NewStyle().
		Bold(true).
		Foreground(color1).
		Border(lipgloss.RoundedBorder(), false, false, true).
		BorderForeground(color1).
		MarginBottom(0)

	color2 := lipgloss.Color(tabColors[1])
	tabStopWatch := lipgloss.NewStyle().
		Bold(true).
		Foreground(color2).
		Border(lipgloss.RoundedBorder(), false, false, true).
		BorderForeground(color2).
		MarginBottom(0)

	header += lipgloss.JoinHorizontal(lipgloss.Center, tabTimer.Render("     ⌛Timer     "), "   ", tabStopWatch.Render("  ⏱  Stop Watch   "))

	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, header, view),
	)
}
