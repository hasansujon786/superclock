package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/cmd"
	"github.com/hasan/superclock/pkg/logger"
)

// DaemonState is like bubble state.Model
type DaemonState struct {
	mu       sync.Mutex
	Timeout  time.Duration // total countdown time
	Interval time.Duration // tick interval
	Elapsed  time.Duration // time passed
	Running  bool
}

func (s *DaemonState) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Running {
		return
	}

	s.Elapsed += s.Interval
	if s.Elapsed > s.Timeout {
		s.Elapsed = s.Timeout
		s.Running = false

		// _ = enc.Encode(state) I want notify client from here with Running false falue
		beeep.Notify("Timer Alert", "Timer completed!", "")
	}

	// Example: notify every 10s
	if int(s.Elapsed.Seconds())%10 == 0 && s.Elapsed.Seconds() != 0 {
		beeep.Notify("Timer Alert",
			fmt.Sprintf("Elapsed: %s", s.Elapsed.String()),
			"")
	}
}

// TODO: check time before play
func (s *DaemonState) Start() {
	s.mu.Lock()
	s.Running = true
	s.mu.Unlock()
}
func (s *DaemonState) Stop() {
	s.mu.Lock()
	s.Running = false
	s.mu.Unlock()
}
func (s *DaemonState) Reset() {
	s.mu.Lock()
	s.Running = false
	s.Elapsed = 0
	s.mu.Unlock()
}
func (s *DaemonState) Toggle() {
	s.mu.Lock()
	s.Running = !s.Running
	s.mu.Unlock()
}
func (s *DaemonState) SetTimer(timeout time.Duration, play any) {
	s.mu.Lock()
	s.Elapsed = 0
	s.Timeout = timeout

	switch v := play.(type) {
	case bool:
		s.Running = v
	}
	s.mu.Unlock()
}

func main() {
	closeLogger := cmd.SetupDotEnv()
	defer closeLogger()

	cmd.RegisterGob()

	state := &DaemonState{
		// Running:  true,
		Timeout:  0,           // default 1 min
		Interval: time.Second, // tick every 1s
	}

	// // Background tick loop
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

func handleConn(conn net.Conn, state *DaemonState) {
	defer conn.Close()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	var req constants.Request
	if err := dec.Decode(&req); err != nil {
		fmt.Println("Error decoding command:", err)
		return
	}

	switch req.Cmd {
	case constants.CmdPause:
		state.Stop()
	case constants.CmdPlay:
		logger.Info("play")
		state.Start()
	case constants.CmdStop:
		state.Reset()
	case constants.CmdToggle:
		state.Toggle()
	case constants.CmdReset:
		state.Reset()
	case constants.CmdSetTimer:
		if pl, ok := req.Payload.(models.CmdSetTimerPayload); ok {
			state.SetTimer(pl.Timeout, pl.Play)
		}
	case constants.CmdGet:
		// no state change
	default:
		fmt.Println("Unknown command:", req.Cmd)
	}

	state.mu.Lock()
	_ = enc.Encode(state)
	state.mu.Unlock()
}
