package cmd

type globalFlags struct {
	Debug      bool
	ConfigFile string
}

var (
	GlobalFlags = &globalFlags{}
)
