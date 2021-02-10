package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Version    string
	Port       int
	DBName     string
	DBUser     string
	DBPassword string
}

func GetConfig() *Config {
	setupEnvVariables()

	cfg := &Config{
		Version:    viper.GetString("version"),
		Port:       viper.GetInt("port"),
		DBName:     viper.GetString("dbname"),
		DBUser:     viper.GetString("dbuser"),
		DBPassword: viper.GetString("dbpassword"),
	}

	return cfg
}

func setupEnvVariables() {
	viper.SetEnvPrefix("hubble")

	err := viper.BindEnv("version")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("port")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("dbname")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("dbuser")
	if err != nil {
		log.Fatal(err)
	}

	err = viper.BindEnv("dbpassword")
	if err != nil {
		log.Fatal(err)
	}
}
