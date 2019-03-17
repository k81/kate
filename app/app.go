package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/k81/kate/configer"
	"github.com/k81/log"
	"github.com/kardianos/osext"
)

var (
	name         string
	homeDir      string
	logFile      string
	errFile      string
	pidFile      string
	confFile     string
	logFormatter log.Formatter = log.PipeKVFormatter

	mctx = log.WithContext(context.Background(), "module", "app")
)

// nolint:gochecknoinits
func init() {
	var (
		bin string
		err error
	)

	if bin, err = osext.Executable(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: get executable path, reason=%v\n", err)
		os.Exit(1)
	}
	homeDir = path.Dir(path.Dir(bin))
	name = path.Base(bin)
	confFile = path.Join(homeDir, "conf", fmt.Sprint(name, ".conf"))
}

// Setup prepare for the application startup
func Setup(configer configer.Configer) {
	var (
		showVersion bool
		err         error
	)

	flag.BoolVar(&showVersion, "version", false, "show version")
	flag.StringVar(&confFile, "config", confFile, "config file name")

	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

	if err = configer.Load(confFile); err != nil {
		log.Fatal(mctx, "load config failed", "error", err)
	}

	pidFile = configer.MustGet("main.pid_file", path.Join(homeDir, "run", fmt.Sprint(name, ".pid")))
	updatePIDFile(pidFile)

	logFile = configer.MustGet("log.log_file", path.Join(homeDir, "log", fmt.Sprint(name, ".log")))
	errFile = configer.MustGet("log.err_file", path.Join(homeDir, "log", fmt.Sprint(name, ".log.wf")))
	initLogger(logFile, errFile, logFormatter)

	log.SetLevelByName(configer.MustGet("log.level", "DEBUG"))

	logVersion()
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

// Cleanup do the application clean up
func Cleanup() {
	// nolint:errcheck
	os.Remove(pidFile)
}

// GetName return the application name
func GetName() string {
	return name
}

// GetHomeDir return the application home directory
func GetHomeDir() string {
	return homeDir
}

// GetConfigFile return the config file used
func GetConfigFile() string {
	return confFile
}

// GetLogFile return the log file path
func GetLogFile() string {
	return logFile
}

// SetLogFile set the log file to f
func SetLogFile(f string) {
	logFile = f
}

// GetErrFile return the error log file path
func GetErrFile() string {
	return errFile
}

// SetErrFile set the error log file to f
func SetErrFile(f string) {
	errFile = f
}

// GetPidFile return the pid file path
func GetPidFile() string {
	return pidFile
}

// SetPidFile set the pid file to f
func SetPidFile(f string) {
	pidFile = f
}

// SetLogFormatter set the logging formatter to f
func SetLogFormatter(f log.Formatter) {
	logFormatter = f
}
