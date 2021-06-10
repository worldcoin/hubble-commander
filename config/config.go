package config

import (
	"log"
	"path"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/spf13/viper"
)

func setupViper() {
	viper.SetConfigFile(getConfigPath())
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Printf("Configuration file not found (%s). Continuing with default config (possibly overridden by env vars).", getConfigPath())
		} else {
			log.Fatalf("failed to read in config: %s", err)
		}
	}
}

func GetConfig() *Config {
	setupViper()

	return &Config{
		Rollup: &RollupConfig{
			SyncSize:               getUint32("rollup.sync_size", 50),
			FeeReceiverPubKeyID:    getUint32("rollup.fee_receiver_pub_key_id", 0),
			TxsPerCommitment:       getUint32("rollup.txs_per_commitment", 32),
			MinCommitmentsPerBatch: getUint32("rollup.min_commitments_per_batch", 1),
			MaxCommitmentsPerBatch: getUint32("rollup.max_commitments_per_batch", 32),
			CommitmentLoopInterval: getDuration("rollup.commitment_loop_interval", 500*time.Millisecond),
			BatchLoopInterval:      getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			GenesisAccounts:        getGenesisAccounts(),
			BootstrapNodeURL:       getStringOrNil("rollup.bootstrap_node_url"),
		},
		API: &APIConfig{
			Version: "0.0.1",
			Port:    getString("api.port", "8080"),
			DevMode: false,
		},
		Postgres: &PostgresConfig{
			Host:           getStringOrNil("postgres.host"),
			Port:           getStringOrNil("postgres.port"),
			Name:           getString("postgres.name", "hubble"),
			User:           getStringOrNil("postgres.user"),
			Password:       getStringOrNil("postgres.password"),
			MigrationsPath: getMigrationsPath(),
		},
		Badger: &BadgerConfig{
			Path: getString("badger.path", getBadgerPath()),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetTestConfig() *Config {
	setupViper()

	return &Config{
		Rollup: &RollupConfig{
			SyncSize:               getUint32("rollup.sync_size", 50),
			FeeReceiverPubKeyID:    getUint32("rollup.fee_receiver_pub_key_id", 0),
			TxsPerCommitment:       getUint32("rollup.txs_per_commitment", 2),
			MinCommitmentsPerBatch: getUint32("rollup.min_commitments_per_batch", 1),
			MaxCommitmentsPerBatch: getUint32("rollup.max_commitments_per_batch", 32),
			CommitmentLoopInterval: getDuration("rollup.commitment_loop_interval", 500*time.Millisecond),
			BatchLoopInterval:      getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			GenesisAccounts:        getGenesisAccounts(),
			BootstrapNodeURL:       getStringOrNil("rollup.bootstrap_node_url"),
		},
		API: &APIConfig{
			Version: "dev-0.0.1",
			Port:    getString("api.port", "8080"),
			DevMode: true,
		},
		Postgres: &PostgresConfig{
			Host:           getStringOrNil("postgres.host"),
			Port:           getStringOrNil("postgres.port"),
			Name:           getString("postgres.name", "hubble_test"),
			User:           getStringOrNil("postgres.user"),
			Password:       getStringOrNil("postgres.password"),
			MigrationsPath: getMigrationsPath(),
		},
		Badger: &BadgerConfig{
			Path: getString("badger.path", getBadgerPath()),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "config.yaml")
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getString("rollup.genesis_path", getGenesisPath())
	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Fatalf("error reading genesis file: %s", err.Error())
	}
	return genesisAccounts
}

func getGenesisPath() string {
	return path.Join(utils.GetProjectRoot(), "genesis.yaml")
}

func getMigrationsPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "postgres", "migrations")
}

func getBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "badger", "data")
}

func getEthereumConfig() *EthereumConfig {
	rpcURL := getStringOrNil("ethereum.rpc_url")
	if rpcURL == nil {
		return nil
	}
	return &EthereumConfig{
		RPCURL:     *rpcURL,
		ChainID:    getStringOrThrow("ethereum.chain_id"),
		PrivateKey: getStringOrThrow("ethereum.private_key"),
	}
}
