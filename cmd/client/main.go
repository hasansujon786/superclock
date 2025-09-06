package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app"
	"github.com/hasan/superclock/app/utils"
	"github.com/hasan/superclock/cmd"
)

func UNUSED(x ...any) {}

func main() {
	closeLogger := cmd.SetupDotEnv()
	defer closeLogger()

	cmd.RegisterGob()

	app := app.NewApp(app.AppViewPomodoro, utils.InitDaemonState())
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("App run failed: %v", err)
		os.Exit(1)
	}
}
