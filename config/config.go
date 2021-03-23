package config

import (
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
	_ = godotenv.Load(".env")
}

func GetConfig() Config {
	return Config{
		Rollup: RollupConfig{
			FeeReceiverIndex: 0,
		},
		API: APIConfig{
			Version: "dev-0.1.0",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
		},
		DB: DBConfig{
			Host:     getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:     getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:     *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble")),
			User:     getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password: getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
		},
		Ethereum: &EthereumConfig{
			RPCURL:     getEnv("ETHEREUM_RPC_URL"),
			ChainID:    getEnv("ETHEREUM_CHAIN_ID"),
			PrivateKey: getEnv("ETHEREUM_PRIVATE_KEY"),
		},
	}
}

func GetTestConfig() Config {
	return Config{
		Rollup: RollupConfig{
			FeeReceiverIndex: 0,
		},
		API: APIConfig{
			Version: "dev-0.1.0",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
		},
		DB: DBConfig{
			Host:     getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:     getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:     *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble_test")),
			User:     getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password: getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
		},
	}
}
