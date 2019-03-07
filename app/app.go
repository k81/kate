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
	"github.com/k81/kate/profiling"
	"github.com/k81/kate/utils"
	"github.com/k81/log"
)

var (
	name         string
	homeDir      string
	logFile      string
	errFile      string
	pidFile      string
	confFile     string
	logFormatter = log.PipeKVFormatter

	logger *log.Logger
	mctx   = context.Background()
)

func init() {
	name = utils.GetAppName()
	homeDir = utils.GetHomeDir()
	logFile = path.Join(homeDir, "log", fmt.Sprint(name, ".log"))
	errFile = path.Join(homeDir, "log", fmt.Sprint(name, ".err"))
	pidFile = path.Join(homeDir, "run", fmt.Sprint(name, ".pid"))
	confFile = path.Join(homeDir, "conf", fmt.Sprint(name, ".yaml"))
}

func Setup() {
	var (
		showVersion = flag.Bool("version", false, "show version")
		confFile    = flag.String("config", confFile, "config file name")
		logAppender *log.FileAppender
		errAppender *log.FileAppender
		err         error
	)

	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if logAppender, err = log.NewFileAppender(log.LevelMask, logFile, logFormatter); err != nil {
		os.Exit(1)
	}
	logAppender.DisableLock()

	if errAppender, err = log.NewFileAppender(log.LevelError|log.LevelFatal, errFile, logFormatter); err != nil {
		os.Exit(1)
	}
	errAppender.DisableLock()

	log.SetLogger(log.NewLogger(logAppender, errAppender))

	logger = log.With("module", app)

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
		syscall.SIGUSR2,
		syscall.SIGHUP,
	)

	for {
		sig := <-sigCh
		logger.Info(mctx, "got signal", "signal", sig)

		switch sig {
		case syscall.SIGINT:
			return sig
		case syscall.SIGQUIT:
			{
				logger.Info(mctx, "fetch done, ready to exit")
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
