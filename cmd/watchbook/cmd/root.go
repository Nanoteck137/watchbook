package cmd

import (
	"log/slog"

	"github.com/nanoteck137/pyrin/trail"
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/config"
	"github.com/spf13/cobra"
)

var logger = trail.NewLogger(&trail.Options{Debug: true, Level: slog.LevelInfo})

var rootCmd = &cobra.Command{
	Use:     watchbook.AppName,
	Version: watchbook.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(watchbook.VersionTemplate(watchbook.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
