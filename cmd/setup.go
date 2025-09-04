package cmd

import (
	"encoding/gob"
	"log"
	"time"

	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/pkg/logger"
	"github.com/joho/godotenv"
)

// Register types that will pass through connection
func RegisterGob() {
	gob.Register(time.Second)
	gob.Register(models.CmdSetTimerPayload{})
}

func SetupDotEnv() func() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing...")
	}

	// Initialize logger
	logger.Init()
	return logger.Close
}
