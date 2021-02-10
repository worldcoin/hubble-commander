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
	var version, dbname, dbuser, dbpassword string
	var port int

	setupEnvVariables()
	readAndSetEnvVariables(
		&version,
		&port,
		&dbname,
		&dbuser,
		&dbpassword,
	)

	cfg := &Config{
		Version:    version,
		Port:       port,
		DBName:     dbname,
		DBUser:     dbuser,
		DBPassword: dbpassword,
	}

	return cfg
}

func CreateConfig(
	version string,
	port int,
	dbname,
	dbuser,
	dbpassword string,
) *Config {
	setupEnvVariables()

	cfg := &Config{
		Version:    version,
		Port:       port,
		DBName:     dbname,
		DBUser:     dbuser,
		DBPassword: dbpassword,
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

func readAndSetEnvVariables(
	version *string,
	port *int,
	dbname,
	dbuser,
	dbpassword *string,
) {
	envVersion := viper.GetString("version")
	if envVersion != "" {
		*version = envVersion
	}

	envPort := viper.GetInt("port")
	if envPort != 0 {
		*port = envPort
	}

	envDBName := viper.GetString("dbname")
	if envDBName != "" {
		*dbname = envDBName
	}


	envDBUser := viper.GetString("dbuser")
	if envDBUser != "" {
		*dbuser = envDBUser
	}

	envDBPassword := viper.GetString("dbpassword")
	if envDBPassword != "" {
		*dbpassword = envDBPassword
	}
}
