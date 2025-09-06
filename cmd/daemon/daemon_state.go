package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hasan/superclock/app"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/pkg/logger"
)

// DaemonStateMutex is like bubble state.Model
type DaemonStateMutex struct {
	mu sync.Mutex
	app.DaemonStateMsg
}

func (s *DaemonStateMutex) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Running {
		return
	}

	s.Elapsed += s.Interval
	if s.Elapsed >= s.Timeout {
		s.Elapsed = s.Timeout
		s.Running = false

		_ = beeep.Notify("Timer Alert", "Timer completed!", "")
	}

	// Example: notify every 10s
	// if int(s.Elapsed.Seconds())%10 == 0 && s.Elapsed.Seconds() != 0 {
	// 	beeep.Notify("Timer Alert", fmt.Sprintf("Elapsed: %s", s.Elapsed.String()), "")
	// }
}

// TODO: check time before play
func (s *DaemonStateMutex) Start() {
	s.mu.Lock()
	s.Running = true
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Stop() {
	s.mu.Lock()
	s.Running = false
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Reset() {
	s.mu.Lock()
	s.Running = false
	s.Elapsed = 0
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Toggle() {
	s.mu.Lock()
	s.Running = !s.Running
	s.mu.Unlock()
}
func (s *DaemonStateMutex) SetTimer(timeout time.Duration, play any) {
	s.mu.Lock()
	s.Elapsed = 0
	s.Timeout = timeout

	switch v := play.(type) {
	case bool:
		s.Running = v
	}
	s.mu.Unlock()
}

func handleConn(conn net.Conn, state *DaemonStateMutex) {

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
	_ = enc.Encode(state.DaemonStateMsg)
	state.mu.Unlock()
}
