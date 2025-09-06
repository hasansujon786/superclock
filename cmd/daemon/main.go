package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/hasan/superclock/app"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/cmd"
)

func main() {
	closeLogger := cmd.SetupDotEnv()
	defer closeLogger()
	cmd.RegisterGob()

	state := &DaemonStateMutex{
		DaemonStateMsg: app.DaemonStateMsg{
			Timeout:  0,           // default timeout
			Interval: time.Second, // tick every 1s
		},
	}

	// Background tick loop
	go func() {
		ticker := time.NewTicker(state.Interval)
		defer ticker.Stop()
		for range ticker.C {
			state.Tick()
		}
	}()

	// TCP server
	ln, err := net.Listen("tcp", constants.Address)
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Daemon running...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}
		go handleConn(conn, state)
	}
}
