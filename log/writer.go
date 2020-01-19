package log

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	RotateSignal = syscall.SIGUSR1
	OpenFlag     = os.O_CREATE | os.O_APPEND | os.O_WRONLY | syscall.O_DSYNC
	OpenPerm     = os.FileMode(0644)
)

type Writer struct {
	location string
	file     *os.File
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWriter(location string) (*Writer, error) {
	file, err := os.OpenFile(location, OpenFlag, OpenPerm)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := &Writer{
		location: location,
		file:     file,
		ctx:      ctx,
		cancel:   cancel,
	}

	w.wg.Add(1)
	go w.startRotateListener()

	return w, nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *Writer) Sync() error {
	return w.file.Sync()
}

func (w *Writer) Close() error {
	w.cancel()
	w.wg.Wait()
	return w.file.Close()
}

func (w *Writer) startRotateListener() {
	defer func() {
		w.wg.Done()
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, RotateSignal)

	for {
		select {
		case <-w.ctx.Done():
			return
		case sig := <-ch:
			if sig == RotateSignal {
				w.rotate()
			}
		}
	}
}

func (w *Writer) rotate() {
	file, err := os.OpenFile(w.location, OpenFlag, OpenPerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate log file \"%s\", reopen, reason=%v", w.location, err)
		return
	}

	if err = syscall.Dup2(int(file.Fd()), int(w.file.Fd())); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to rotate log file \"%s\", dup2(), reason=%v", w.location, err)
		return
	}

	if err = file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to close log file \"%s\", close(), reason=%v", w.location, err)
		return
	}
}
