package models

import "time"

type CmdSetTimerPayload struct {
	Timeout time.Duration
	Play    any // bool | nil
}
