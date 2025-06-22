package cmd

import (
	"log"

	"github.com/nanoteck137/watchbook"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     watchbook.CliAppName,
	Version: watchbook.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(watchbook.VersionTemplate(watchbook.CliAppName))
}
