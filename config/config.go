package config

type Config struct {
	Version    string
	Port       int
	DBName     string
	DBUser     string
	DBPassword string
}

func GetConfig() *Config {
	cfg := &Config{
		Version:    "dev-0.1.0",
		Port:       8080,
		DBName:     "hubble_test",
		DBUser:     "hubble",
		DBPassword: "root",
	}

	return cfg
}
