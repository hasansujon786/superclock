package stopwatch

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/ui"
	"github.com/hasan/superclock/app/utils"
)

type StopWatchModel struct {
	stopwatch stopwatch.Model
	paused    bool
	laps      []lap

	width, height  int
	CurrentAppView constants.AppView
}

func NewStopWatchModel() StopWatchModel {
	return StopWatchModel{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

func (m StopWatchModel) Init() tea.Cmd {
	return nil
}

func (m StopWatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case stopwatch.TickMsg:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		return m, cmd

	case stopwatch.StartStopMsg:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		if !m.stopwatch.Running() && !utils.DurationEnded(m.stopwatch.Elapsed()) {
			m.paused = true
		} else if m.paused {
			m.paused = false
		}
		return m, cmd

	case stopwatch.ResetMsg:
		var cmd tea.Cmd
		m.stopwatch, cmd = m.stopwatch.Update(msg)
		m.paused = false
		return m, cmd

	case tea.KeyMsg:
		if constants.AppViewStopWatch != m.CurrentAppView {
			return m, nil
		}

		switch msg.String() {

		case " ":
			if !m.stopwatch.Running() {
				return m, nil
			}

			nLap := createNewLap(m.stopwatch.Elapsed(), m.laps)
			m.laps = append([]lap{nLap}, m.laps...)
			return m, nil
		case "r":
			m.laps = []lap{}
			return m, tea.Batch(m.stopwatch.Stop(), m.stopwatch.Reset())
		case "s":
			return m, m.stopwatch.Toggle()
		}
	}

	return m, nil
}

// ------------------ View Function ----------------

func (m StopWatchModel) View() string {
	cWidth := styles.ContainerStyle.GetWidth()
	cHeight := styles.ContainerStyle.GetHeight()

	clkState := constants.ClockState{Running: m.stopwatch.Running(), Paused: m.paused}

	timeDigit := ui.TimerDigit(
		utils.FormatStopwatch(m.stopwatch.Elapsed()),
		cWidth,
		constants.NerdFont,
	)

	ctlBox := buildControlBox(&clkState)

	header := styles.HeaderStyle.Render("  ⏱  Stop Watch   ")
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		timeDigit,
		" ",
		buildLapItemList(m.laps),
		" ",
		ctlBox,
	)

	viewBox := styles.ContainerStyle.Render(
		lipgloss.Place(cWidth, cHeight, lipgloss.Center, lipgloss.Center, content),
	)

	help := buildHelp(m.stopwatch.Running())

	return lipgloss.JoinVertical(lipgloss.Center, header, viewBox, help)
}

func buildHelp(running bool) string {

	s := utils.If(running, "s: stop ", "s: start")
	s += " • r: reset\nspace: lap • q: quit"

	return styles.FooterStyle.Render(s)
}

// ------------------ Lap Functions -------------------

type lap struct {
	time  time.Duration
	diff  time.Duration
	index int
}

func buildLapItem(item lap) string {
	greyStyle := lipgloss.NewStyle().Foreground(styles.ThemeColors.Muted)
	greenStyle := lipgloss.NewStyle().Foreground(styles.ThemeColors.Success)

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		utils.SpaceBetween(
			30,
			fmt.Sprintf("%v   %02d ", greenStyle.Render(""), item.index),
			greyStyle.Render("+"+utils.FormatStopwatch(item.diff)),
			utils.FormatStopwatch(item.time),
		)...,
	)
}

func buildLapItemList(laps []lap) string {
	if len(laps) == 0 {
		return ""
	}

	lapsComp := []string{}
	for _, lap := range laps {
		lapsComp = append(lapsComp, buildLapItem(lap))
	}
	return lipgloss.JoinVertical(lipgloss.Center, lapsComp...)
}

func buildControlBox(clState *constants.ClockState) string {
	playPauseBtn := ui.ButtonStyles.Render(utils.If(clState.IsRunning(), "   ", "   "))

	if clState.IsStopped() {
		return playPauseBtn
	}

	lapStopBtn := ui.ButtonStyles.Render(utils.If(clState.IsPaused(), "   ", "   "))
	return lipgloss.JoinHorizontal(lipgloss.Center, playPauseBtn, "  ", lapStopBtn)
}

func createNewLap(lapTime time.Duration, laps []lap) lap {
	index := len(laps) + 1
	var diff time.Duration // only here time.Duration is not a type [NotAType]
	if index == 1 {
		diff = lapTime
	} else {
		diff = lapTime - laps[0].time
	}

	return lap{time: lapTime, index: index, diff: diff}
}
