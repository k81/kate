package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

type StdLogger struct {
	sync.Mutex
	formatter Formatter
}

func NewStdLogger(formatter Formatter) *StdLogger {
	return &StdLogger{
		formatter: formatter,
	}
}

func (l *StdLogger) Log(entry *Entry) {
	var (
		buf []byte
		err error
	)

	l.Lock()
	defer l.Unlock()

	if buf, err = l.formatter.Format(entry); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to format entry: reason=[%v]\n", err)
		return
	}

	if _, err = io.Copy(os.Stderr, bytes.NewBuffer(buf)); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to copy data to writer: reason=[%v]\n", err)
		return
	}
}
