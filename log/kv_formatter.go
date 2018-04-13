package log

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"github.com/modern-go/gls"
)

var (
	PID         int
	KVFormatter = &kvFormatter{}
)

func init() {
	PID = os.Getpid()
}

type kvFormatter struct{}

func (l *kvFormatter) Format(entry *Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "%-7s %s [%d][%d] ",
		entry.Level.String(),
		entry.Time.Format("2006-01-02 15:04:05.000"),
		PID,
		gls.GoID(),
	)

	b.WriteString("msg=[")
	b.WriteString(entry.Msg)
	b.WriteString("] ")

	for i := 0; i < len(entry.KeyVals); i += 2 {
		key, val := entry.KeyVals[i], entry.KeyVals[i+1]
		b.WriteString(toString(key))
		b.WriteString("=[")
		fmt.Fprint(b, val)
		b.WriteString("] ")
	}
	b.WriteString("fileline=[")
	b.WriteString(fmt.Sprint(path.Base(entry.File), ":", entry.Line))
	b.WriteString("]")
	b.WriteByte('\n')

	return b.Bytes(), nil
}
