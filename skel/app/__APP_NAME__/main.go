package main

import (
	"context"
	"fmt"
	"os"

	"github.com/k81/kate/app"
	"github.com/k81/log"
	"github.com/spf13/cobra"

	"__PROJECT_DIR__/app/__APP_NAME__/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use:        app.GetName(),
		Short:      "Sample Application",
		SuggestFor: []string{app.GetName()},
	}

	mctx = log.WithContext(context.Background(), "module", "main")
)

func main() {
	cobra.EnablePrefixMatching = true
	rootCmd.PersistentFlags().BoolVar(&cmd.GlobalFlags.Debug, "debug", false, "enable debug")
	rootCmd.PersistentFlags().StringVar(&cmd.GlobalFlags.ConfigFile, "config", app.GetDefaultConfigFile(), "config file path")

	rootCmd.AddCommand(
		cmd.NewVersionCmd(),
		cmd.NewStartCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
	}
}
