package config

import (
	"os"

	"github.com/spf13/viper"
)

func getEnvOrDefault(key string, fallback *string) *string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return &value
}

func getFromViperOrDefault(key string, fallback *string) *string {
	value := viper.GetString(key)
	if value == "" {
		return fallback
	}
	return &value
}
