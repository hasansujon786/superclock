package models

import (
	"time"

	"github.com/hasan/superclock/app/constants"
)

type PomodoroStateMsg struct {
	TotalTime time.Duration // total countdown time
	Timeout   time.Duration // How long until the timer expires.
	ModeIdx   int           // Mode name index
	Running   bool
}

type DaemonStateMsg struct {
	Interval time.Duration // tick interval
	Pomodoro PomodoroStateMsg
}

type Request struct {
	Cmd     constants.Command
	Payload any
}

type CmdSetTimerPayload struct {
	ModeIdx int
	Play    any // bool | nil
	Timeout time.Duration
}
