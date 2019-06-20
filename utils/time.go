package utils

import (
	"fmt"
	"time"
)

// Milliseconds return the milliseconds of time
func Milliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// TimeLocationOfUTCOffset return the time.Location of utc offset
func TimeLocationOfUTCOffset(utcOffset int) *time.Location {
	zoneName := fmt.Sprintf("UTC%+d", utcOffset)
	return time.FixedZone(zoneName, utcOffset*60*60)
}

// TimeInUTCOffset return time in specified utc offset
func TimeInUTCOffset(t time.Time, utcOffset int) time.Time {
	return t.In(TimeLocationOfUTCOffset(utcOffset))
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

// GetMonthsOfDayRange return the month string of day range [beginDay, endDay]
func GetMonthsOfDayRange(layout string, beginDay, endDay time.Time) []string {
	monthMap := make(map[string]bool)
	for curDay := beginDay; !curDay.After(endDay); curDay = curDay.Add(24 * time.Hour) {
		month := curDay.Format(layout)
		monthMap[month] = true
	}

	result := make([]string, 0, len(monthMap))
	for month := range monthMap {
		result = append(result, month)
	}
	return result
}
