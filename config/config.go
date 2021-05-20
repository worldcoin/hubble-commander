package config

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("CI") != "true" {
		loadDotEnv()
	}
}

func loadDotEnv() {
	dotEnvFilePath := path.Join(utils.GetProjectRoot(), ".env")
	_ = godotenv.Load(dotEnvFilePath)
}

func getMigrationsPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "postgres", "migrations")
}

func getBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "badger", "data")
}

func getGenesisPath() string {
	return path.Join(utils.GetProjectRoot(), "genesis.yaml")
}

func GetConfig() Config {
	return Config{
		Rollup: RollupConfig{
			SyncBatches:             true,
			FeeReceiverPubKeyID:     0,
			TxsPerCommitment:        32,
			MinCommitmentsPerBatch:  1,
			MaxCommitmentsPerBatch:  32,
			CommitmentLoopInterval:  500 * time.Millisecond,
			BatchLoopInterval:       500 * time.Millisecond,
			BlockNumberLoopInterval: 500 * time.Millisecond,
			GenesisAccounts:         getGenesisAccounts(),
		},
		API: APIConfig{
			Version: "0.0.1",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: false,
		},
		Postgres: PostgresConfig{
			Host:           getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:           getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:           *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble")),
			User:           getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password:       getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
			MigrationsPath: *getEnvOrDefault("HUBBLE_MIGRATIONS_PATH", ref.String(getMigrationsPath())),
		},
		Badger: BadgerConfig{
			Path: *getEnvOrDefault("HUBBLE_BADGER_PATH", ref.String(getBadgerPath())),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetTestConfig() Config {
	return Config{
		Rollup: RollupConfig{
			SyncBatches:             false,
			FeeReceiverPubKeyID:     0,
			TxsPerCommitment:        2,
			MinCommitmentsPerBatch:  1,
			MaxCommitmentsPerBatch:  32,
			CommitmentLoopInterval:  500 * time.Millisecond,
			BatchLoopInterval:       500 * time.Millisecond,
			BlockNumberLoopInterval: 500 * time.Millisecond,
			GenesisAccounts:         getGenesisAccounts(),
		},
		API: APIConfig{
			Version: "dev-0.0.1",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: true,
		},
		Postgres: PostgresConfig{
			Host:           getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:           getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:           *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble_test")),
			User:           getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password:       getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
			MigrationsPath: getMigrationsPath(),
		},
		Badger: BadgerConfig{
			Path: *getEnvOrDefault("HUBBLE_BADGER_PATH", ref.String(getBadgerPath())),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := *getEnvOrDefault("HUBBLE_GENESIS_PATH", ref.String(getGenesisPath()))

	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Fatalf("error reading genesis file: %s", err.Error())
	}

	return genesisAccounts
}

func getEthereumConfig() *EthereumConfig {
	rpcURL := getEnvOrDefault("ETHEREUM_RPC_URL", nil)
	if rpcURL == nil {
		return nil
	}
	return &EthereumConfig{
		RPCURL:     *rpcURL,
		ChainID:    getEnv("ETHEREUM_CHAIN_ID"),
		PrivateKey: getEnv("ETHEREUM_PRIVATE_KEY"),
	}
}
