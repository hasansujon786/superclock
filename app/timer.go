package app

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Timer model ---
type TimerModel struct {
	value int
	count int
}

func (m TimerModel) Init() tea.Cmd {
	return nil
}

func (m TimerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "k": // Exit app
			m.count++
			return m, nil

		case "j": // Exit app
			m.count--
			return m, nil
		}

	}

	return m, nil
}

func (m TimerModel) View() string {
	s := "⏱ Timer running...\n\n"

	s += strconv.Itoa(m.count)

	return s
}

