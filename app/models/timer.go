package models

import "time"

// ------------------------------------------------
// -- PickerValue ---------------------------------
// ------------------------------------------------
type PickerValue struct {
	Hour    int
	Minute  int
	Second int
}

// ToDuration converts PickerValue to time.Duration
func (p PickerValue) ToDuration() time.Duration {
	return time.Duration(p.Hour)*time.Hour +
		time.Duration(p.Minute)*time.Minute +
		time.Duration(p.Second)*time.Second
}

func (p PickerValue) IsEmpty() bool {
	return p.Second == 0 && p.Minute == 0 && p.Hour == 0
}
