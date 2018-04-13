package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/kardianos/osext"
)

var (
	homeDir string
	appName string
)

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
	appName = path.Base(bin)
}

func GetHomeDir() string {
	return homeDir
}

func GetAppName() string {
	return appName
}
