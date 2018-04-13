package app

import (
	"fmt"

	"github.com/k81/kate/log"
)

var (
	VER_MAJOR   string
	VER_MINOR   string
	VER_PATCH   string
	BUILD_DATE  string
	REVISION    string
	LAST_AUTHOR string
	LAST_DATE   string
)

func printVersion() {
	fmt.Println("Version:    ", fmt.Sprintf("%s.%s.%s", VER_MAJOR, VER_MINOR, VER_PATCH))
	fmt.Println("Revision:   ", REVISION)
	fmt.Println("Last Author:", LAST_AUTHOR)
	fmt.Println("Last Date:  ", LAST_DATE)
	fmt.Println("Build Date: ", BUILD_DATE)
}

func logVersion() {
	log.Info(mctx, "app info",
		"version", fmt.Sprintf("%s.%s.%s", VER_MAJOR, VER_MINOR, VER_PATCH),
		"revision", REVISION,
		"last_author", LAST_AUTHOR,
		"last_date", LAST_DATE,
		"build_date", BUILD_DATE,
	)
}
