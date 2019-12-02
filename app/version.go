package app

import (
	"fmt"

	"go.uber.org/zap"
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

func PrintVersion() {
	fmt.Println("Version:    ", fmt.Sprintf("%s.%s.%s", VersionMajor, VersionMinor, VersionPatch))
	fmt.Println("Revision:   ", Revision)
	fmt.Println("Last Author:", LastAuthor)
	fmt.Println("Last Date:  ", LastDate)
	fmt.Println("Build Date: ", BuildDate)
}

func LogVersion(logger *zap.Logger) {
	logger.Info("app info",
		zap.String("version", fmt.Sprintf("%s.%s.%s", VersionMajor, VersionMinor, VersionPatch)),
		zap.String("revision", Revision),
		zap.String("last_author", LastAuthor),
		zap.String("last_date", LastDate),
		zap.String("build_date", BuildDate),
	)
}
