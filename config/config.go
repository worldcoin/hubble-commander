package config

import (
	"log"
	"os"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("CI") != "true" {
		loadDotEnv()
	}
}

func loadDotEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Config struct {
	Version    string
	Port       string
	DBHost     *string
	DBPort     *string
	DBName     string
	DBUser     *string
	DBPassword *string
}

func GetConfig() Config {
	return Config{
		Version:    "dev-0.1.0",
		Port:       *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
		DBHost:     getEnvOrDefault("HUBBLE_DBHOST", nil),
		DBPort:     getEnvOrDefault("HUBBLE_DBPORT", nil),
		DBName:     *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble")),
		DBUser:     getEnvOrDefault("HUBBLE_DBUSER", nil),
		DBPassword: getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
	}
}

func GetTestConfig() Config {
	return Config{
		Version:    "dev-0.1.0",
		Port:       *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
		DBHost:     getEnvOrDefault("HUBBLE_DBHOST", nil),
		DBPort:     getEnvOrDefault("HUBBLE_DBPORT", nil),
		DBName:     *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble_test")),
		DBUser:     getEnvOrDefault("HUBBLE_DBUSER", nil),
		DBPassword: getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
	}
}

func getEnvOrDefault(name string, fallback *string) *string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	return &value
}
