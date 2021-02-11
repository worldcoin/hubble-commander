package config

import (
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Version    string
	Port       string
	DBHost     *string
	DBPort     *string
	DBName     *string
	DBUser     *string
	DBPassword *string
}

func GetConfig() *Config {
	viper.SetEnvPrefix("hubble")

	cfg := &Config{
		Version:    "dev-0.1.0",
		Port:       *getEnvOrDefault("port", utils.MakeStringPointer("8080")),
		DBHost:     getEnvOrDefault("dbhost", utils.MakeStringPointer("localhost")),
		DBPort:     getEnvOrDefault("dbport", nil),
		DBName:     getEnvOrDefault("dbname", utils.MakeStringPointer("hubble")),
		DBUser:     getEnvOrDefault("dbuser", nil),
		DBPassword: getEnvOrDefault("dbpassword", nil),
	}

	return cfg
}

func GetTestConfig() *Config {
	viper.SetEnvPrefix("hubble")

	cfg := &Config{
		Version:    "dev-0.1.0",
		Port:       *getEnvOrDefault("port", utils.MakeStringPointer("8080")),
		DBHost:     getEnvOrDefault("dbhost", utils.MakeStringPointer("localhost")),
		DBPort:     getEnvOrDefault("dbport", nil),
		DBName:     getEnvOrDefault("dbname", utils.MakeStringPointer("hubble_test")),
		DBUser:     getEnvOrDefault("dbuser", nil),
		DBPassword: getEnvOrDefault("dbpassword", nil),
	}

	return cfg
}

func getEnvOrDefault(name string, def *string) *string {
	err := viper.BindEnv(name)
	if err != nil {
		return def
	}

	val := viper.GetString(name)
	if val == "" {
		return def
	}

	return utils.MakeStringPointer(val)
}
