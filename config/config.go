package config

import (
	"os"
	"strconv"
)

type Config struct {
	Version    string
	Port       int
	DBName     string
	DBUser     string
	DBPassword string
}

func GetConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("Port"))
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Version:    os.Getenv("Version"),
		Port:       port,
		DBName:     os.Getenv("DBName"),
		DBUser:     os.Getenv("DBUser"),
		DBPassword: os.Getenv("DBPassword"),
	}

	return cfg, err
}
