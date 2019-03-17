package app

import (
	"fmt"

	"github.com/k81/log"
)

var (
	// VersionMajor the major version number
	VersionMajor string
	// VersionMinor the minor version number
	VersionMinor string
	// VersionPatch the patch version number
	VersionPatch string
	// BuildDate the build date
	BuildDate string
	// Revision the revision number
	Revision string
	// LastAuthor the author of last commit
	LastAuthor string
	// LastDate the date of last commit
	LastDate string
)

func printVersion() {
	fmt.Println("Version:    ", fmt.Sprintf("%s.%s.%s", VersionMajor, VersionMinor, VersionPatch))
	fmt.Println("Revision:   ", Revision)
	fmt.Println("Last Author:", LastAuthor)
	fmt.Println("Last Date:  ", LastDate)
	fmt.Println("Build Date: ", BuildDate)
}

func logVersion() {
	log.Info(mctx, "app info",
		"version", fmt.Sprintf("%s.%s.%s", VersionMajor, VersionMinor, VersionPatch),
		"revision", Revision,
		"last_author", LastAuthor,
		"last_date", LastDate,
		"build_date", BuildDate,
	)
}
