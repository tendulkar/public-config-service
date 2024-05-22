// config/load.go
package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func LoadConfig(logger *slog.Logger) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file", slog.Any("error", err))
		os.Exit(1)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		logger.Error("Error unmarshaling config", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("Application Config loaded", slog.Any("config", config))

	return &config
}
