package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/kardianos/osext"
)

var (
	name     string
	homeDir  string
	pidFile  string
	confFile string
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
	confFile = path.Join(homeDir, "conf", fmt.Sprint(name, ".ini"))
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
func UpdatePIDFile(fileName string) error {
	var (
		runDir = path.Dir(fileName)
		pid    = os.Getpid()
		err    error
	)

	if err = os.MkdirAll(runDir, 0755); err != nil {
		return fmt.Errorf("failed to create dir: dir=%v, error=%w", runDir, err)
	}

	if err = ioutil.WriteFile(fileName, []byte(strconv.Itoa(pid)), 0666); err != nil {
		return fmt.Errorf("failed to write pid: file=%v, pid=%v, error=%v", fileName, pid, err)
	}

	pidFile = fileName

	return nil
}

// RemovePIDFile do the application clean up
func RemovePIDFile() {
	if pidFile != "" {
		// nolint:errcheck
		_ = os.Remove(pidFile)
	}
}
