package app

import (
	"fmt"
	"os"
	"path"

	"github.com/k81/log"
)

func updatePIDFile(pidFile string) {
	var (
		runDir = path.Dir(pidFile)
		pid    = os.Getpid()
		f      *os.File
		err    error
	)

	if _, err = os.Stat(runDir); err != nil && os.IsNotExist(err) {
		if err = os.Mkdir(runDir, 0755); err != nil {
			log.Error(mctx, "mkdir dir", "dir", runDir, "error", err)
			return
		}
	}

	f, err = os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Error(mctx, "write pid file", "file", pidFile, "error", err)
		return
	}
	fmt.Fprint(f, pid)
	// nolint:errcheck
	f.Close()
	log.Info(mctx, "pid file written", "file", pidFile, "pid", pid)
}
