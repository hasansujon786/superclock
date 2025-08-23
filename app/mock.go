package app

import "time"

var tt = time.Now()
var DVEFAUTL_LAPS = []lap{
	{index: 03, diff: time.Since(tt), time: time.Since(tt)},
	{index: 02, diff: time.Since(tt), time: time.Since(tt)},
	{index: 01, diff: time.Since(tt), time: time.Since(tt)},
}
