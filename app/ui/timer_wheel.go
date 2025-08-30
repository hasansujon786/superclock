package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/styles"
)

type CursorPositon int

const (
	CursorPosHour CursorPositon = iota
	CursorPosMinute
	CursorPosSecond
	CursorPosNone
)

type TimerWheelModel struct {
	Position CursorPositon
	Value    PickerValue
}

func NewTimerWheelModel(focus CursorPositon) TimerWheelModel {
	return TimerWheelModel{
		Position: focus,
		Value:    PickerValue{hour: 0, minute: 0, seconsd: 0},
	}
}

func (tp *TimerWheelModel) PickerMoveCursorLeft() {
	if tp.Position == CursorPosNone {
		return
	}
	tp.Position = (tp.Position - 1 + CursorPosNone) % CursorPosNone
}
func (tp *TimerWheelModel) PickerMoveCursorRight() {
	if tp.Position == CursorPosNone {
		return
	}
	tp.Position = (tp.Position + 1) % CursorPosNone
}

func (tp *TimerWheelModel) IncreaseValue() {
	if tp.Position == CursorPosNone {
		return
	}

	switch tp.Position {
	case CursorPosHour:
		if tp.Value.hour == 23 {
			tp.Value.hour = 0
		} else {
			tp.Value.hour++
		}
	case CursorPosMinute:
		if tp.Value.minute == 59 {
			tp.Value.minute = 0
		} else {
			tp.Value.minute++
		}
	case CursorPosSecond:
		if tp.Value.seconsd == 59 {
			tp.Value.seconsd = 0
		} else {
			tp.Value.seconsd++
		}
	}
}
func (tp *TimerWheelModel) DecreaseValue() {
	if tp.Position == CursorPosNone {
		return
	}

	switch tp.Position {
	case CursorPosHour:
		if tp.Value.hour == 0 {
			tp.Value.hour = 23
		} else {
			tp.Value.hour--
		}
	case CursorPosMinute:
		if tp.Value.minute == 0 {
			tp.Value.minute = 59
		} else {
			tp.Value.minute--
		}
	case CursorPosSecond:
		if tp.Value.seconsd == 0 {
			tp.Value.seconsd = 59
		} else {
			tp.Value.seconsd--
		}

	}
}

func (tp *TimerWheelModel) Focus(pos CursorPositon) {
	tp.Position = pos
}
func (tp *TimerWheelModel) FocusLast() {
	if tp.Position == CursorPosNone {
		tp.Focus(CursorPosSecond)
	} else {
		tp.Focus(tp.Position)
	}
}
func (tp *TimerWheelModel) Blur() {
	tp.Position = CursorPosNone
}

func (tp *TimerWheelModel) Reset() {
	tp.Value = PickerValue{hour: 0, minute: 0, seconsd: 0}
}
func (tp *TimerWheelModel) ResetCurrent() {
	if tp.Position == CursorPosNone {
		return
	}

	switch tp.Position {
	case CursorPosHour:
		tp.Value.hour = 0
	case CursorPosMinute:
		tp.Value.minute = 0
	case CursorPosSecond:
		tp.Value.seconsd = 0
	}
}

// ------------------------------------------------
// -- PickerValue ---------------------------------
// ------------------------------------------------
type PickerValue struct {
	hour    int
	minute  int
	seconsd int
}

// ToDuration converts PickerValue to time.Duration
func (p PickerValue) ToDuration() time.Duration {
	return time.Duration(p.hour)*time.Hour +
		time.Duration(p.minute)*time.Minute +
		time.Duration(p.seconsd)*time.Second
}

// ------------------------------------------------
// -- TimerWhell components -----------------------
// ------------------------------------------------
func TimerWhell(timer PickerValue, cursor CursorPositon) string {
	Unfocused := lipgloss.NewStyle().
		Foreground(styles.ThemeColors.Primary)

	focused := lipgloss.NewStyle().
		Foreground(styles.ThemeColors.Secondary).
		BorderForeground(styles.ThemeColors.Secondary).
		Underline(true)

	hourDigit := Unfocused.Render(fmt.Sprintf("%02d", timer.hour))
	minuteDigit := Unfocused.Render(fmt.Sprintf("%02d", timer.minute))
	secondDigit := Unfocused.Render(fmt.Sprintf("%02d", timer.seconsd))
	separator := Unfocused.Render(":")

	switch cursor {
	case CursorPosHour:
		hourDigit = focused.Render(fmt.Sprintf("%02d", timer.hour))
	case CursorPosMinute:
		minuteDigit = focused.Render(fmt.Sprintf("%02d", timer.minute))
	case CursorPosSecond:
		secondDigit = focused.Render(fmt.Sprintf("%02d", timer.seconsd))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, hourDigit, separator, minuteDigit, separator, secondDigit)
}
