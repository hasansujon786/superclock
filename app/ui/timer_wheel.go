package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/hasan/superclock/app/models"
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
	Position     CursorPositon
	LastPosition CursorPositon
	Value        models.PickerValue
}

func NewTimerWheelModel(focus CursorPositon) TimerWheelModel {
	return TimerWheelModel{
		Position:     focus,
		LastPosition: focus,
		Value:        models.PickerValue{Hour: 0, Minute: 0, Second: 0},
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
		if tp.Value.Hour == 23 {
			tp.Value.Hour = 0
		} else {
			tp.Value.Hour++
		}
	case CursorPosMinute:
		if tp.Value.Minute == 59 {
			tp.Value.Minute = 0
		} else {
			tp.Value.Minute++
		}
	case CursorPosSecond:
		if tp.Value.Second == 59 {
			tp.Value.Second = 0
		} else {
			tp.Value.Second++
		}
	}
}
func (tp *TimerWheelModel) DecreaseValue() {
	if tp.Position == CursorPosNone {
		return
	}

	switch tp.Position {
	case CursorPosHour:
		if tp.Value.Hour == 0 {
			tp.Value.Hour = 23
		} else {
			tp.Value.Hour--
		}
	case CursorPosMinute:
		if tp.Value.Minute == 0 {
			tp.Value.Minute = 59
		} else {
			tp.Value.Minute--
		}
	case CursorPosSecond:
		if tp.Value.Second == 0 {
			tp.Value.Second = 59
		} else {
			tp.Value.Second--
		}

	}
}

func (tp *TimerWheelModel) Focus(pos CursorPositon) {
	tp.Position = pos
}
func (tp *TimerWheelModel) FocusLast() {
	if tp.LastPosition == CursorPosNone {
		tp.Focus(CursorPosSecond)
	} else {
		tp.Focus(tp.LastPosition)
	}
}
func (tp *TimerWheelModel) Blur() {
	tp.LastPosition = tp.Position
	tp.Position = CursorPosNone
}

func (tp *TimerWheelModel) Reset() {
	tp.Value = models.PickerValue{Hour: 0, Minute: 0, Second: 0}
}
func (tp *TimerWheelModel) ResetCurrent() {
	if tp.Position == CursorPosNone {
		return
	}

	switch tp.Position {
	case CursorPosHour:
		tp.Value.Hour = 0
	case CursorPosMinute:
		tp.Value.Minute = 0
	case CursorPosSecond:
		tp.Value.Second = 0
	}
}

// ------------------------------------------------
// -- TimerWhell components -----------------------
// ------------------------------------------------
func TimerWhell(timer models.PickerValue, cursor CursorPositon) string {
	unfocused := lipgloss.NewStyle().
		Foreground(styles.ThemeColors.Primary)

	focused := lipgloss.NewStyle().
		Foreground(styles.ThemeColors.Secondary).
		BorderForeground(styles.ThemeColors.Secondary).
		Underline(true)

	hourDigit := unfocused.Render(fmt.Sprintf("%02d", timer.Hour))
	minuteDigit := unfocused.Render(fmt.Sprintf("%02d", timer.Minute))
	secondDigit := unfocused.Render(fmt.Sprintf("%02d", timer.Second))
	separator := unfocused.Render(":")

	switch cursor {
	case CursorPosHour:
		hourDigit = focused.Render(fmt.Sprintf("%02d", timer.Hour))
	case CursorPosMinute:
		minuteDigit = focused.Render(fmt.Sprintf("%02d", timer.Minute))
	case CursorPosSecond:
		secondDigit = focused.Render(fmt.Sprintf("%02d", timer.Second))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, hourDigit, separator, minuteDigit, separator, secondDigit)
}
