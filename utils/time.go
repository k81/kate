package utils

import (
	"fmt"
	"time"
)

// Milliseconds return the milliseconds of time
func Milliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// TimeInUTCOffset return time in specified utc offset
func TimeInUTCOffset(t time.Time, utcOffset int) time.Time {
	zoneName := fmt.Sprintf("UTC%+d", utcOffset)
	loc := time.FixedZone(zoneName, utcOffset*60*60)
	return t.In(loc)
}

// GetDayRangeOfMonth [firstDay, lastDay]
func GetDayRangeOfMonth(date time.Time) (firstDay, lastDay time.Time) {
	year, month, _ := date.Date()
	firstDay = time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	lastDay = firstDay.AddDate(0, 1, -1)
	return firstDay, lastDay
}

// GetTimeRangeOfDay [begin, end)
func GetTimeRangeOfDay(t time.Time) (begin, end time.Time) {
	year, month, day := t.Date()
	begin = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	end = begin.AddDate(0, 0, 1)
	return begin, end
}
