package config

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getString(key, fallback string) string {
	viper.SetDefault(key, fallback)
	return viper.GetString(key)
}

func getStringOrNil(key string) *string {
	value := viper.GetString(key)
	if value == "" {
		return nil
	}
	return &value
}

func getStringOrThrow(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Fatalf("%s config not specified", key)
	}
	return value
}

func getUint32(key string, fallback uint32) uint32 {
	viper.SetDefault(key, fallback)
	return viper.GetUint32(key)
}

func getUint64(key string, fallback uint32) uint64 {
	viper.SetDefault(key, fallback)
	return viper.GetUint64(key)
}

func getBool(key string, fallback bool) bool {
	viper.SetDefault(key, fallback)
	return viper.GetBool(key)
}

func getDuration(key string, fallback time.Duration) time.Duration {
	viper.SetDefault(key, fallback)
	return viper.GetDuration(key)
}
