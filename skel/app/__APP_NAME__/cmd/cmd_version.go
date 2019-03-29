package cmd

import (
	"github.com/k81/kate/app"
	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print app version",
		Run:   versionCmdFunc,
	}
	return cmd
}

func versionCmdFunc(cmd *cobra.Command, args []string) {
	app.PrintVersion()
}
