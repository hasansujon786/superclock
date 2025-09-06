package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/hasan/superclock/app/constants"
	"github.com/hasan/superclock/app/models"
	"github.com/hasan/superclock/pkg/logger"
)

// DaemonStateMutex is like bubble state.Model
type DaemonStateMutex struct {
	mu sync.Mutex
	models.DaemonStateMsg
}

func (s *DaemonStateMutex) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Pomodoro.Running {
		s.Pomodoro.Elapsed += s.Interval
		if s.Pomodoro.Elapsed >= s.Pomodoro.Timeout {
			s.Pomodoro.Elapsed = s.Pomodoro.Timeout
			s.Pomodoro.Running = false

			_ = beeep.Notify("Timer Alert", "Timer completed!", "")
		}
	}

	// Example: notify every 10s
	// if int(s.Elapsed.Seconds())%10 == 0 && s.Elapsed.Seconds() != 0 {
	// 	beeep.Notify("Timer Alert", fmt.Sprintf("Elapsed: %s", s.Elapsed.String()), "")
	// }
}

// TODO: check time before play
func (s *DaemonStateMutex) Start() {
	s.mu.Lock()
	s.Pomodoro.Running = true
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Stop() {
	s.mu.Lock()
	s.Pomodoro.Running = false
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Reset() {
	s.mu.Lock()
	s.Pomodoro.Running = false
	s.Pomodoro.Elapsed = 0
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Toggle() {
	s.mu.Lock()
	s.Pomodoro.Running = !s.Pomodoro.Running
	s.mu.Unlock()
}
func (s *DaemonStateMutex) SetTimer(timeout time.Duration, play any) {
	s.mu.Lock()
	s.Pomodoro.Elapsed = 0
	s.Pomodoro.Timeout = timeout

	switch v := play.(type) {
	case bool:
		s.Pomodoro.Running = v
	}
	s.mu.Unlock()
}

func handleConn(conn net.Conn, state *DaemonStateMutex) {

	defer conn.Close()

	dec := gob.NewDecoder(conn)
	enc := gob.NewEncoder(conn)

	var req models.Request
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
