package config

import (
	"os"

	"github.com/nanoteck137/watchbook"
	"github.com/spf13/viper"
)

var logger = watchbook.DefaultLogger()

type Config struct {
	LibraryDir string `mapstructure:"library_dir"`
}

func setDefaults() {
	viper.BindEnv("library_dir")
}

func validateConfig(config *Config) {
	hasError := false

	validate := func(expr bool, msg string) {
		if expr {
			logger.Error("Config Validation", "err", msg)
			hasError = true
		}
	}

	// NOTE(patrik): Has default value, here for completeness
	// validate(config.RunMigrations == "", "run_migrations needs to be set")
	validate(config.LibraryDir == "", "library_dir needs to be set")

	if hasError {
		os.Exit(1)
	}
}

var ConfigFile string
var LoadedConfig Config

func InitConfig() {
	setDefaults()

	if ConfigFile != "" {
		viper.SetConfigFile(ConfigFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix(watchbook.CliAppName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Warn("failed to load config", "err", err)
	}

	err = viper.Unmarshal(&LoadedConfig)
	if err != nil {
		logger.Fatal("failed to unmarshal config", "err", err)
	}

	configCopy := LoadedConfig
	logger.Debug("Current Config", "config", configCopy)

	validateConfig(&LoadedConfig)
}
