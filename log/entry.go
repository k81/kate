package log

import (
	"errors"
	"time"
)

var ErrMissingValue = errors.New("(MISSING)")

type Entry struct {
	Level   Level
	Time    time.Time
	File    string
	Line    int
	Msg     string
	KeyVals []interface{}
}
