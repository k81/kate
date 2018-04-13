package log

import (
	"bytes"
	"fmt"
	"path"

	"github.com/modern-go/gls"
)

var (
	DidiFormatter = &didiFormatter{}
)

const DidiTimeLayout = "2006-01-02T15:04:05.999+07:00"

type didiFormatter struct{}

func (l *didiFormatter) Format(entry *Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "[%s][%s][%d]",
		entry.Level.String(),
		entry.Time.Format(DidiTimeLayout),
		gls.GoID(),
	)

	b.WriteString("msg=")
	b.WriteString(entry.Msg)
	b.WriteString("||")

	for i := 0; i < len(entry.KeyVals); i += 2 {
		key, val := entry.KeyVals[i], entry.KeyVals[i+1]
		b.WriteString(toString(key))
		b.WriteString("=")
		fmt.Fprint(b, val)
		b.WriteString("||")
	}
	b.WriteString("fileline=")
	b.WriteString(fmt.Sprint(path.Base(entry.File), ":", entry.Line))
	b.WriteByte('\n')

	return b.Bytes(), nil
}
