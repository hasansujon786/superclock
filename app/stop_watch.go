package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	boxStyle := lipgloss.NewStyle().
		Width(40).
		Height(12).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1)

	contentWidth := boxStyle.GetWidth() - 2

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		// Background(lipgloss.Color("240")).
		Padding(0, 1).
		MarginBottom(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))

	time := lipgloss.Place(
		contentWidth, 2,
		lipgloss.Center, lipgloss.Center,
		titleStyle.Render(formatStopwatch(m.stopwatch.Elapsed())),
	)

	laps := []string{}
	for _, lap := range m.laps {
		laps = append(laps, buildLapItem(lap))
	}

	lapsList := lipgloss.JoinVertical(lipgloss.Center, laps...)

	content := lipgloss.JoinVertical(lipgloss.Center, time, lapsList)
	box := boxStyle.Render(content)
	help := buildHelp(m.stopwatch.Running())

	return lipgloss.JoinVertical(lipgloss.Center, box, help)
}

func buildHelp(running bool) string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	s := If(running, "s: stop ", "s: start")
	s += " • r: reset • space: lap • q: quit"

	return footerStyle.Render(s)
}

// ------------------ Lap Function -------------------

type lap struct {
	time  time.Duration
	diff  time.Duration
	index int
}

func buildLapItem(item lap) string {
	greyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		SpaceBetween(
			30,
			fmt.Sprintf("%v   %02d ", greyStyle.Render("▸"), item.index),
			greyStyle.Render("+"+formatStopwatch(item.diff)),
			formatStopwatch(item.time),
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
