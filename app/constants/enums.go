package constants

// ------------------------------------------------
// -- Clock State ---------------------------------
// ------------------------------------------------
type ClockState struct {
	Running bool
	Paused  bool
}

var (
	ClockStateStoped  = ClockState{Running: false, Paused: false}
	ClockStateRunning = ClockState{Running: true, Paused: false}
	ClockStatePaused  = ClockState{Running: false, Paused: true}
)

func (state ClockState) IsRunning() bool {
	return state == ClockStateRunning
}
func (state ClockState) IsStopped() bool {
	return state == ClockStateStoped
}
func (state ClockState) IsPaused() bool {
	return state == ClockStatePaused
}

// ------------------------------------------------
// -- Daemon Actions ------------------------------
// ------------------------------------------------
type Command string

const (
	CmdGet      Command = "get"
	CmdSetTimer Command = "set_timer"
	CmdPlay     Command = "play"
	CmdPause    Command = "pause"
	CmdStop     Command = "stop"
	CmdToggle   Command = "toggle"
	CmdReset    Command = "reset"
)

type Request struct {
	Cmd     Command
	Payload any
}
