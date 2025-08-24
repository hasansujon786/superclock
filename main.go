package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/clock-tui/app"
)

func main() {
	app := app.NewApp()
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("App run failed: %v", err)
		os.Exit(1)
	}
}
