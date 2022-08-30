package config

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cast"
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

func getStringOrPanic(key string) string {
	value := viper.GetString(key)
	if value == "" {
		log.Panicf("%s config not specified", key)
	}
	return value
}

func getAddressOrNil(key string) *common.Address {
	strAddress := getStringOrNil(key)
	if strAddress == nil {
		return nil
	}
	address := common.HexToAddress(*strAddress)
	return &address
}

func getUint32(key string, fallback uint32) uint32 {
	viper.SetDefault(key, fallback)
	return viper.GetUint32(key)
}

func getUint64(key string, fallback uint64) uint64 {
	viper.SetDefault(key, fallback)
	return viper.GetUint64(key)
}

func getUint64OrPanic(key string) uint64 {
	value, err := cast.ToUint64E(viper.Get(key))
	if err != nil {
		log.Panicf("failed to read %s config: %v", key, err)
	}
	return value
}

//nolint:unparam
func getBool(key string, fallback bool) bool {
	viper.SetDefault(key, fallback)
	return viper.GetBool(key)
}

func getDuration(key string, fallback time.Duration) time.Duration {
	viper.SetDefault(key, fallback)
	return viper.GetDuration(key)
}
