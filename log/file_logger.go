package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	OpenFlag = os.O_CREATE | os.O_APPEND | os.O_WRONLY | syscall.O_DSYNC
	OpenPerm = 0644
)

type FileLogger struct {
	sync.Mutex
	disableLock bool // enable log line bigger than 4k
	logFileName string
	errFileName string
	logFile     *os.File
	errFile     *os.File
	formatter   Formatter
}

func NewFileLogger(logFileName, errFileName string, formatter Formatter) (l *FileLogger, err error) {
	l = &FileLogger{
		logFileName: logFileName,
		errFileName: errFileName,
		formatter:   formatter,
	}
	if l.logFile, err = os.OpenFile(logFileName, OpenFlag, OpenPerm); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to open log file: reason=[%v]\n", err)
		return
	}
	if l.errFile, err = os.OpenFile(errFileName, OpenFlag, OpenPerm); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to open err file: reason=[%v]\n", err)
	}
	return
}

func (l *FileLogger) DisableLock() {
	l.disableLock = true
}

func (l *FileLogger) Log(entry *Entry) {
	var (
		buf []byte
		err error
	)

	if buf, err = l.formatter.Format(entry); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to format entry: reason=[%v]\n", err)
		return
	}

	if !l.disableLock {
		l.Lock()
		defer l.Unlock()
	}

	if entry.Level >= ErrorLevel {
		if _, err = io.Copy(l.errFile, bytes.NewBuffer(buf)); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: failed to copy data to err file writer: reason=[%v]\n", err)
		}
	}

	if _, err = io.Copy(l.logFile, bytes.NewBuffer(buf)); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to copy data to log file writer: reason=[%v]\n", err)
	}
}

func (l *FileLogger) Rotate() {
	var (
		logFile *os.File
		errFile *os.File
		err     error
	)

	if !l.disableLock {
		l.Lock()
		defer l.Unlock()
	}

	if errFile, err = os.OpenFile(l.errFileName, OpenFlag, OpenPerm); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate err file \"%s\", reopen, reason=%v", l.errFileName, err)
		return
	}

	if err = syscall.Dup2(int(errFile.Fd()), int(l.errFile.Fd())); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate err file \"%s\", dup2(), reason=%v", l.errFileName, err)
		return
	}

	if err = errFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to close err file \"%s\", close(), reason=%v", l.errFileName, err)
		return
	}

	if logFile, err = os.OpenFile(l.logFileName, OpenFlag, OpenPerm); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate log file\"%s\", reopen, reason=%v", l.logFileName, err)
		return
	}

	if err = syscall.Dup2(int(logFile.Fd()), int(l.logFile.Fd())); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate log file \"%s\", dup2(), reason=%v", l.logFileName, err)
		return
	}

	if err = logFile.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to close log file \"%s\", close(), reason=%v", l.logFileName, err)
		return
	}
}

func (l *FileLogger) SetRotateSignal(signo os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signo)

	go func() {
		for {
			sig := <-ch

			switch {
			case sig == syscall.SIGUSR1:
				Info(mctx, "got signal, rotating", "signal", signo)
				l.Rotate()
				Info(mctx, "log file rotated")
			}
		}
	}()
}
