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

// DaemonState is like bubble timer.Model
type DaemonState struct {
	mu       sync.Mutex
	Timeout  time.Duration // total countdown time
	Interval time.Duration // tick interval
	Elapsed  time.Duration // time passed
	Running  bool
}

func (t *DaemonState) Tick() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.Running {
		return
	}

	t.Elapsed += t.Interval
	if t.Elapsed > t.Timeout {
		t.Elapsed = t.Timeout
		t.Running = false
		beeep.Notify("Timer Alert", "Timer completed!", "")
	}

	// Example: notify every 10s
	if int(t.Elapsed.Seconds())%10 == 0 && t.Elapsed.Seconds() != 0 {
		beeep.Notify("Timer Alert",
			fmt.Sprintf("Elapsed: %s", t.Elapsed.String()),
			"")
	}
}

// TODO: check time before play
func (t *DaemonState) Start() {
	t.mu.Lock()
	t.Running = true
	t.mu.Unlock()
}
func (t *DaemonState) Stop() {
	t.mu.Lock()
	t.Running = false
	t.mu.Unlock()
}
func (t *DaemonState) Reset() {
	t.mu.Lock()
	t.Running = false
	t.Elapsed = 0
	t.mu.Unlock()
}
func (t *DaemonState) Toggle() {
	t.mu.Lock()
	t.Running = !t.Running
	t.mu.Unlock()
}
func (t *DaemonState) SetTimer(timeout time.Duration) {
	t.mu.Lock()
	t.Timeout = timeout
	t.mu.Unlock()
}

func main() {
	closeLogger := cmd.SetupDotEnv()
	defer closeLogger()

	cmd.RegisterGob()

	timer := &DaemonState{
		Timeout:  time.Minute, // default 1 min
		Interval: time.Second, // tick every 1s
	}

	// Background tick loop
	go func() {
		ticker := time.NewTicker(timer.Interval)
		defer ticker.Stop()
		for range ticker.C {
			timer.Tick()
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
		go handleConn(conn, timer)
	}
}

func handleConn(conn net.Conn, timer *DaemonState) {
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
		timer.Stop()
	case constants.CmdPlay:
		timer.Start()
	case constants.CmdStop:
		timer.Reset()
	case constants.CmdToggle:
		timer.Toggle()
	case constants.CmdReset:
		timer.Reset()
	case constants.CmdSetTimer:
		logger.Info("pickerValue")
		if pickerValue, ok := req.Payload.(models.PickerValue); ok {
			timer.SetTimer(pickerValue.ToDuration())
		}
	case constants.CmdGet:
		// no state change
	default:
		fmt.Println("Unknown command:", req.Cmd)
	}

	timer.mu.Lock()
	_ = enc.Encode(timer)
	timer.mu.Unlock()
}
