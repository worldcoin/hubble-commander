package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s env var not specified", key)
	}
	return value
}

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
