package config

import (
	"log"
	"os"
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
