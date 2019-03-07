package utils

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/k81/log"
)

func UpdatePIDFile(ctx context.Context, pidFile string) {
	var (
		runDir string = path.Dir(pidFile)
		f      *os.File
		pid    int = os.Getpid()
		err    error
	)

	if _, err = os.Stat(runDir); os.IsNotExist(err) {
		if err = os.Mkdir(runDir, 0755); err != nil {
			log.Error(ctx, "mkdir dir", "dir", runDir, "error", err)
			return
		}
	}

	f, err = os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(ctx, "write pid file", "file", pidFile, "error", err)
		return
	}
	fmt.Fprint(f, pid)
	f.Close()
	log.Info(ctx, "pid file written", "file", pidFile, "pid", pid)
}
