package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type lap struct {
	time  time.Duration
	diff  time.Duration
	index int
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

// --- StopWatch model ---
type StopWatchModel struct {
	stopwatch     stopwatch.Model
	laps          []lap
	width, height int
}

func NewStopWatch() StopWatchModel {
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

func (m StopWatchModel) View() string {
	s := "⏳ Stopwatch ticking...\n\n"

	s += "       "
	s += formatStopwatch(m.stopwatch.Elapsed()) + "\n\n"

	for _, lap := range m.laps {

		s += fmt.Sprintf("%02d    +%v    %v\n\n", lap.index, formatStopwatch(lap.diff), formatStopwatch(lap.time))
	}

	if m.stopwatch.Running() {
		s += "s: stop  "
	} else {
		s += "s: start "
	}
	s += "| r: reset | q: quit\n"
	return s
}

func formatStopwatch(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	ms := int(d.Milliseconds()/10) % 100 // hundredths of a second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d.%02d", h, m, s, ms)
	}
	return fmt.Sprintf("%02d:%02d.%02d", m, s, ms)
}
