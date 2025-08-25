package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/components"
	"github.com/hasan/superclock/app/styles"
	"github.com/hasan/superclock/app/utils"
)

type StopWatchModel struct {
	stopwatch     stopwatch.Model
	laps          []lap
	width, height int
}

func DefaultStopWatchModel() StopWatchModel {
	return StopWatchModel{
		stopwatch: stopwatch.NewWithInterval(time.Millisecond),
	}
}

func (m StopWatchModel) Init() tea.Cmd {
	return nil
}

func (m StopWatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
			return m, m.stopwatch.Reset()
		case "s":
			return m, m.stopwatch.Toggle()
		}
	}

	var cmd tea.Cmd
	m.stopwatch, cmd = m.stopwatch.Update(msg)
	return m, cmd
}

// ------------------ View Function ----------------

func (m StopWatchModel) View() string {
	contentWidth := styles.ContainerStyle.GetWidth() - 2
	timeDigit := components.TimerDigit(utils.FormatStopwatch(m.stopwatch.Elapsed()), contentWidth, components.NerdFont)

	laps := []string{}
	for _, lap := range m.laps {
		laps = append(laps, buildLapItem(lap))
	}
	lapList := lipgloss.JoinVertical(lipgloss.Center, laps...)

	header := styles.HeaderStyle.Render("  ⏱  Stop Watch   ")
	content := lipgloss.JoinVertical(lipgloss.Center, timeDigit, " ", lapList)
	box := styles.ContainerStyle.Render(content)
	help := buildHelp(m.stopwatch.Running())

	return lipgloss.JoinVertical(lipgloss.Center, header, box, help)
}

func buildHelp(running bool) string {

	s := utils.If(running, "s: stop ", "s: start")
	s += " • r: reset • space: lap • q: quit"

	return styles.FooterStyle.Render(s)
}

// ------------------ Lap Function -------------------

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
