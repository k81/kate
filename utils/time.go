package utils

import "time"

// Milliseconds return the milliseconds of time
func Milliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
