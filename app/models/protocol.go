package models

import (
	"time"

	"github.com/hasan/superclock/app/constants"
)

type PomodoroStateMsg struct {
	Timeout time.Duration // total countdown time
	Elapsed time.Duration // time passed
	Running bool
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
	Timeout time.Duration
	Play    any // bool | nil
}
