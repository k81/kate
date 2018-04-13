package log

import (
	"runtime"
	"time"
)

type LogContext struct {
	keyvals   []interface{}
	hasValuer bool
}

func NewLogContext() *LogContext {
	return &LogContext{}
}

func (c *LogContext) With(keyvals ...interface{}) *LogContext {
	if len(keyvals) == 0 {
		return c
	}
	kvs := append(c.keyvals, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	return &LogContext{
		keyvals:   kvs[:len(kvs):len(kvs)],
		hasValuer: c.hasValuer || containsValuer(keyvals),
	}
}

func (c *LogContext) WithPrefix(keyvals ...interface{}) *LogContext {
	if len(keyvals) == 0 {
		return c
	}
	n := len(c.keyvals) + len(keyvals)
	if len(keyvals)%2 != 0 {
		n++
	}
	kvs := make([]interface{}, 0, n)
	kvs = append(kvs, keyvals...)
	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}
	kvs = append(kvs, c.keyvals...)
	return &LogContext{
		keyvals:   kvs,
		hasValuer: c.hasValuer || containsValuer(keyvals),
	}
}

func (c *LogContext) newEntry(lvl Level, msg string, keyvals []interface{}) (entry *Entry) {
	var (
		kvs []interface{}
		ok  bool
	)

	kvs = make([]interface{}, 0, len(c.keyvals)+len(keyvals)+2)
	kvs = append(kvs, c.keyvals...)
	kvs = append(kvs, keyvals...)

	if len(kvs)%2 != 0 {
		kvs = append(kvs, ErrMissingValue)
	}

	if c.hasValuer {
		bindValues(kvs[:len(c.keyvals)])
	}

	entry = &Entry{
		Time:    time.Now(),
		Level:   lvl,
		KeyVals: kvs,
		Msg:     msg,
	}

	_, entry.File, entry.Line, ok = runtime.Caller(2)
	if !ok {
		entry.File = "???"
		entry.Line = -1
	}

	return entry
}
