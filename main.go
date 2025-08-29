package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hasan/superclock/app"
	"github.com/hasan/superclock/logger"
	"github.com/joho/godotenv"
)

func UNUSED(x ...any) {}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing...")
	}

	// Initialize logger
	logger.Init()
	defer logger.Close()

	app := app.NewApp(app.AppViewStopWatch)
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Printf("App run failed: %v", err)
		os.Exit(1)
	}
}
