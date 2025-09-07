package config

import (
	"os"

	"github.com/nanoteck137/watchbook"
	"github.com/nanoteck137/watchbook/types"
	"github.com/spf13/viper"
)

var logger = watchbook.DefaultLogger()

type Config struct {
	RunMigrations   bool   `mapstructure:"run_migrations"`
	ListenAddr      string `mapstructure:"listen_addr"`
	DataDir         string `mapstructure:"data_dir"`
	Username        string `mapstructure:"username"`
	InitialPassword string `mapstructure:"initial_password"`
	JwtSecret       string `mapstructure:"jwt_secret"`
}

func (c *Config) WorkDir() types.WorkDir {
	return types.WorkDir(c.DataDir)
}

func setDefaults() {
	viper.SetDefault("run_migrations", "true")
	viper.SetDefault("listen_addr", ":3000")
	viper.BindEnv("data_dir")
	viper.BindEnv("username")
	viper.BindEnv("initial_password")
	viper.BindEnv("jwt_secret")
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
	validate(config.ListenAddr == "", "listen_addr needs to be set")
	validate(config.DataDir == "", "data_dir needs to be set")
	validate(config.Username == "", "username needs to be set")
	validate(config.InitialPassword == "", "initial_password needs to be set")
	validate(config.JwtSecret == "", "jwt_secret needs to be set")

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

	viper.SetEnvPrefix(watchbook.AppName)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Warn("failed to load config", "err", err)
	}

	err = viper.Unmarshal(&LoadedConfig)
	if err != nil {
		logger.Fatal("failed to unmarshal config", "err", err)
	}

	hide := func(s string) string {
		var res string
		for i := 0; i < len(s); i++ {
			res += "*"
		}

		return res
	}

	configCopy := LoadedConfig
	configCopy.JwtSecret = hide(configCopy.JwtSecret)
	configCopy.InitialPassword = hide(configCopy.InitialPassword)

	logger.Debug("Current Config", "config", configCopy)

	validateConfig(&LoadedConfig)
}
