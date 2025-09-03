package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app/features/timer"
	"github.com/hasan/superclock/cmd"
)

func main() {
	closeLogger := cmd.SetupDotEnv()
	defer closeLogger()

	cmd.RegisterGob()

	p := tea.NewProgram(timer.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error starting client:", err)
		os.Exit(1)
	}
}
