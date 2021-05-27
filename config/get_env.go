package config

import (
	"github.com/spf13/viper"
)

func getFromViperOrDefault(key string, fallback *string) *string {
	value := viper.GetString(key)
	if value == "" {
		return fallback
	}
	return &value
}
