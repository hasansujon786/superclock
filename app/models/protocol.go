package models

import (
	"time"

	"github.com/hasan/superclock/app/constants"
)

type DaemonStateMsg struct {
	Timeout  time.Duration // total countdown time
	Interval time.Duration // tick interval
	Elapsed  time.Duration // time passed
	Running  bool
}

type Request struct {
	Cmd     constants.Command
	Payload any
}

type CmdSetTimerPayload struct {
	Timeout time.Duration
	Play    any // bool | nil
}
