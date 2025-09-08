package pomodoro

import "time"

type pomodoroView int

const (
	viewDashboard pomodoroView = iota
	viewTimer
	viewEdit
)

type pomoPlayMode struct {
	name string
	time time.Duration
}

type TimerResetMsg struct{}
