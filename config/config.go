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
	"github.com/spf13/viper"
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
			SyncBatches:            true,
			FeeReceiverPubKeyID:    0,
			TxsPerCommitment:       32,
			MinCommitmentsPerBatch: 1,
			MaxCommitmentsPerBatch: 32,
			CommitmentLoopInterval: 500 * time.Millisecond,
			BatchLoopInterval:      500 * time.Millisecond,
			GenesisAccounts:        getGenesisAccounts(),
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
			SyncBatches:            false,
			FeeReceiverPubKeyID:    0,
			TxsPerCommitment:       2,
			MinCommitmentsPerBatch: 1,
			MaxCommitmentsPerBatch: 32,
			CommitmentLoopInterval: 500 * time.Millisecond,
			BatchLoopInterval:      500 * time.Millisecond,
			GenesisAccounts:        getGenesisAccounts(),
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
	viper.SetEnvPrefix("ETHEREUM")
	rpcURL := viper.GetString("rpc_url")
	if len(rpcURL) < 1 {
		return nil
	}
	return &EthereumConfig{
		RPCURL:     rpcURL,
		ChainID:    viper.GetString("chain_id"),
		PrivateKey: viper.GetString("private_key"),
	}
}

func getViperConfig() *Config {
	viper.SetConfigFile(path.Join(utils.GetProjectRoot(), "config.yml"))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read in config: %s", err)
	}

	return &Config{
		Rollup: RollupConfig{
			Prune:                   viper.GetBool("prune"),
			SyncBatches:             viper.GetBool("sync_batches"),
			FeeReceiverPubKeyID:     viper.GetUint32("fee_receiver_pub_key_id"),
			TxsPerCommitment:        viper.GetUint32("txs_per_commitment"),
			MinCommitmentsPerBatch:  viper.GetUint32("min_commitments_per_batch"),
			MaxCommitmentsPerBatch:  viper.GetUint32("max_commitments_per_batch"),
			CommitmentLoopInterval:  viper.GetDuration("commitment_loop_interval"),
			BatchLoopInterval:       viper.GetDuration("batch_loop_interval"),
			BlockNumberLoopInterval: viper.GetDuration("block_number_loop_interval"),
			GenesisAccounts:         getGenesisAccounts(),
		},
		API: APIConfig{
			Version: viper.GetString("version"),
			Port:    viper.GetString("port"),
			DevMode: viper.GetBool("dev_mode"),
		},
		Postgres: PostgresConfig{
			Host:           getFromViperOrDefault("dbhost", nil),
			Port:           getFromViperOrDefault("dbport", nil),
			Name:           viper.GetString("dbname"),
			User:           getFromViperOrDefault("dbuser", nil),
			Password:       getFromViperOrDefault("dbpassword", nil),
			MigrationsPath: *getFromViperOrDefault("migrations_path", ref.String(getMigrationsPath())),
		},
		Badger: BadgerConfig{
			Path: *getFromViperOrDefault("badger_path", ref.String(getBadgerPath())),
		},
		Ethereum: getEthereumConfig(),
	}
}
