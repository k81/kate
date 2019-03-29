package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"strconv"
	"syscall"

	"github.com/k81/log"
	"github.com/kardianos/osext"
)

var (
	name     string
	homeDir  string
	pidFile  string
	confFile string

	mctx = log.WithContext(context.Background(), "module", "app")
)

// nolint:gochecknoinits
func init() {
	var (
		bin string
		err error
	)

	if bin, err = osext.Executable(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: get executable path, reason=%v\n", err)
		os.Exit(1)
	}
	homeDir = path.Dir(path.Dir(bin))
	name = path.Base(bin)
	confFile = path.Join(homeDir, "conf", fmt.Sprint(name, ".conf"))
}

// Wait wait for signal to interrupt application running
func Wait() os.Signal {
	sigCh := make(chan os.Signal, 2)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGUSR2,
		syscall.SIGHUP,
	)

	for {
		sig := <-sigCh
		log.Info(mctx, "got signal", "signal", sig)

		switch sig {
		case syscall.SIGINT:
			return sig
		case syscall.SIGQUIT:
			{
				log.Info(mctx, "fetch done, ready to exit")
				return sig
			}
		case syscall.SIGTERM:
			return sig
		case syscall.SIGUSR2:
			continue
		case syscall.SIGHUP:
			continue
		}
	}
}

// GetName return the application name
func GetName() string {
	return name
}

// GetHomeDir return the application home directory
func GetHomeDir() string {
	return homeDir
}

// GetDefaultConfigFile return the config file used
func GetDefaultConfigFile() string {
	return confFile
}

// GetPidFile return the pid file path
func GetPidFile() string {
	return pidFile
}

// UpdatePIDFile update the pid in pidfile
func UpdatePIDFile(fileName string) {
	var (
		runDir = path.Dir(fileName)
		pid    = os.Getpid()
		err    error
	)

	if err = os.MkdirAll(runDir, 0755); err != nil {
		log.Error(mctx, "create run dir", "dir", runDir, "error", err)
		return
	}

	if err = ioutil.WriteFile(fileName, []byte(strconv.Itoa(pid)), 0666); err != nil {
		log.Error(mctx, "write pid to file", "pid", pid, "file", fileName, "error", err)
		return
	}
	pidFile = fileName
	log.Info(mctx, "pid file written", "file", fileName, "pid", pid)
}

// RemovePIDFile do the application clean up
func RemovePIDFile() {
	if pidFile != "" {
		// nolint:errcheck
		_ = os.Remove(pidFile)
	}
}
