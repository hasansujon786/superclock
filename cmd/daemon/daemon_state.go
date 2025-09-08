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
	"github.com/hasan/superclock/app/utils"
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
		logger.Info("Tick timeout", s.Pomodoro.Timeout)

		s.Pomodoro.Timeout -= s.Interval
		if s.Pomodoro.Timeout <= 0 {
			s.Pomodoro.Timeout = 0
			s.Pomodoro.Running = false

			beeep.AppName = "SuperClock"
			_ = beeep.Notify(
				"Pomodoro",
				fmt.Sprintf("Timeout %v", utils.FormatDurationHumanize(s.Pomodoro.TotalTime)),
				"testdata/warning.png",
			)
			// beeep.Notify()
		}
	}

	// Example: notify every 10s
	// if int(s.Elapsed.Seconds())%10 == 0 && s.Elapsed.Seconds() != 0 {
	// beeep.Notify("Timer Alert", fmt.Sprintf("Elapsed: %s", s.Elapsed.String()), "")
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
	s.Pomodoro.Timeout = s.Pomodoro.TotalTime
	s.mu.Unlock()
}
func (s *DaemonStateMutex) Toggle(timeout time.Duration) {
	s.mu.Lock()
	s.Pomodoro.Timeout = timeout
	s.Pomodoro.Running = !s.Pomodoro.Running
	s.mu.Unlock()
}
func (s *DaemonStateMutex) SetTimer(playload any) {
	if pl, ok := playload.(models.CmdSetTimerPayload); ok {
		s.mu.Lock()

		s.Pomodoro.Timeout = pl.Timeout
		s.Pomodoro.TotalTime = pl.Timeout
		s.Pomodoro.ModeIdx = pl.ModeIdx

		switch v := pl.Play.(type) {
		case bool:
			s.Pomodoro.Running = v
		}

		s.mu.Unlock()
	}
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
	case constants.CmdPlay:
		state.Start()
	case constants.CmdStop:
		state.Stop()
	case constants.CmdToggle:
		if timeout, ok := req.Payload.(time.Duration); ok {
			state.Toggle(timeout)
		}
	case constants.CmdReset:
		state.Reset()
	case constants.CmdSetTimer:
		state.SetTimer(req.Payload)
	case constants.CmdGet:
		// no state change
	default:
		fmt.Println("Unknown command:", req.Cmd)
	}

	state.mu.Lock()
	_ = enc.Encode(state.DaemonStateMsg)
	state.mu.Unlock()
}
