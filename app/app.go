package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/k81/kate/config"
	"github.com/k81/kate/log"
	"github.com/k81/kate/profiling"
	"github.com/k81/kate/utils"
)

var (
	name         string
	homeDir      string
	logFile      string
	errFile      string
	pidFile      string
	confFile     string
	logFormatter log.Formatter
	logger       *log.FileLogger

	mctx = log.SetContext(context.Background(), "module", "app")
)

func init() {
	name = utils.GetAppName()
	homeDir = utils.GetHomeDir()
	logFile = path.Join(homeDir, "log", fmt.Sprint(name, ".log"))
	errFile = path.Join(homeDir, "log", fmt.Sprint(name, ".err"))
	pidFile = path.Join(homeDir, "run", fmt.Sprint(name, ".pid"))
	confFile = path.Join(homeDir, "conf", fmt.Sprint(name, ".yaml"))
	logFormatter = log.DidiFormatter
}

func Setup() {
	var (
		showVersion = flag.Bool("version", false, "show version")
		confFile    = flag.String("config", confFile, "config file name")
		err         error
	)

	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if logger, err = log.NewFileLogger(logFile, errFile, logFormatter); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: create log writer, reason=%v", err)
		os.Exit(1)
	}
	logger.DisableLock()
	logger.SetRotateSignal(syscall.SIGUSR1)
	log.SetLogger(logger)

	logVersion()

	config.Init(name, *confFile)
	log.SetLevelByName(GetLogLevel())

	utils.UpdatePIDFile(mctx, pidFile)

	if profiling.Enabled() {
		profiling.Start()
	}
}

func Wait() os.Signal {
	sigCh := make(chan os.Signal, 2)

	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGUSR1,
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
		case syscall.SIGUSR1:
			{
				logger.Rotate()
				log.Info(mctx, "log file rotated")
			}
		case syscall.SIGUSR2:
			continue
		case syscall.SIGHUP:
			continue
		}
	}
}

func Cleanup() {
	os.Remove(pidFile)
}

func GetName() string {
	return name
}

func GetHomeDir() string {
	return homeDir
}

func GetConfigFile() string {
	return confFile
}

func GetLogFile() string {
	return logFile
}

func SetLogFile(f string) {
	logFile = f
}

func GetErrFile() string {
	return errFile
}

func SetErrFile(f string) {
	errFile = f
}

func GetPidFile() string {
	return pidFile
}

func SetPidFile(f string) {
	pidFile = f
}

func SetLogFormatter(f log.Formatter) {
	logFormatter = f
}
