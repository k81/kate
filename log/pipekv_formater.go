package log

import (
	"bytes"
	"fmt"
	"path"

	"github.com/modern-go/gls"
)

var (
	PipeKVFormatter = &pipeKvFormatter{}
)

type pipeKvFormatter struct{}

func (l *pipeKvFormatter) Format(entry *Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	fmt.Fprintf(b, "[%s][%s][%d][%s:%d][tag=%s][",
		entry.Level.String(),
		entry.Time.Format("2006-01-02 15:04:05.000"),
		gls.GoID(),
		path.Base(entry.File),
		entry.Line,
		entry.Msg,
	)

	for i := 0; i < len(entry.KeyVals); i += 2 {
		key, val := entry.KeyVals[i], entry.KeyVals[i+1]
		b.WriteString(toString(key))
		b.WriteString("=")
		fmt.Fprint(b, val)
		b.WriteString("||")
	}

	if len(entry.KeyVals) > 0 {
		b.Truncate(b.Len() - 2)
	}

	b.WriteString("]\n")

	return b.Bytes(), nil
}
