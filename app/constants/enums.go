package constants

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
