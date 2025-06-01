package cmd

import (
	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     watchbook.AppName,
	Version: watchbook.Version,
}

func Execute() {
	logger := watchbook.DefaultLogger()

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to run root command", "err", err)
	}
}

func init() {
	rootCmd.SetVersionTemplate(watchbook.VersionTemplate(watchbook.AppName))

	cobra.OnInitialize(config.InitConfig)

	rootCmd.PersistentFlags().StringVarP(&config.ConfigFile, "config", "c", "", "Config File")
}
