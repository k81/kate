package app

import (
	"os"

	"github.com/k81/log"
)

func initLogger(logFile, errFile string, formatter log.Formatter) {
	logAppender, err := log.NewFileAppender(log.LevelMask, logFile, formatter)
	if err != nil {
		os.Exit(1)
	}
	logAppender.DisableLock()

	errAppender, err := log.NewFileAppender(log.LevelError|log.LevelFatal, errFile, formatter)
	if err != nil {
		os.Exit(1)
	}
	errAppender.DisableLock()

	log.SetLogger(log.NewLogger(logAppender, errAppender))
}
