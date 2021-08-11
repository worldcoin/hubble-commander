package config

import (
	"path"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	log "github.com/sirupsen/logrus"
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
		Log: getLogConfig(),
		Bootstrap: &BootstrapConfig{
			Prune:            false, // overridden in main
			GenesisAccounts:  getGenesisAccounts(),
			BootstrapNodeURL: getStringOrNil("bootstrap.node_url"),
		},
		Rollup: &RollupConfig{
			SyncSize:               getUint32("rollup.sync_size", 50),
			FeeReceiverPubKeyID:    getUint32("rollup.fee_receiver_pub_key_id", 0),
			MinTxsPerCommitment:    getUint32("rollup.min_txs_per_commitment", 1),
			MaxTxsPerCommitment:    getUint32("rollup.max_txs_per_commitment", 32),
			MinCommitmentsPerBatch: getUint32("rollup.min_commitments_per_batch", 1),
			MaxCommitmentsPerBatch: getUint32("rollup.max_commitments_per_batch", 32),
			CommitmentLoopInterval: getDuration("rollup.commitment_loop_interval", 500*time.Millisecond),
			BatchLoopInterval:      getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			DisableSignatures:      false, // overridden in main
		},
		API: &APIConfig{
			Version: "0.0.1",
			Port:    getString("api.port", "8080"),
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
	return &Config{
		Log: &LogConfig{
			Level:  log.InfoLevel,
			Format: "text",
		},
		Bootstrap: &BootstrapConfig{
			Prune:            false,
			GenesisAccounts:  readGenesisAccounts(getGenesisPath()),
			BootstrapNodeURL: nil,
		},
		Rollup: &RollupConfig{
			SyncSize:               50,
			FeeReceiverPubKeyID:    0,
			MinTxsPerCommitment:    2,
			MaxTxsPerCommitment:    2,
			MinCommitmentsPerBatch: 1,
			MaxCommitmentsPerBatch: 32,
			CommitmentLoopInterval: 500 * time.Millisecond,
			BatchLoopInterval:      500 * time.Millisecond,
			DisableSignatures:      true,
		},
		API: &APIConfig{
			Version: "dev-0.0.1",
			Port:    "8080",
		},
		Postgres: &PostgresConfig{
			Host:           nil,
			Port:           nil,
			Name:           "hubble_test",
			User:           ref.String("hubble"),
			Password:       ref.String("root"),
			MigrationsPath: getMigrationsPath(),
		},
		Badger: &BadgerConfig{
			Path: getBadgerPath(),
		},
		Ethereum: nil,
	}
}

func getConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "config.yaml")
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getString("rollup.genesis_path", getGenesisPath())
	return readGenesisAccounts(filename)
}

func readGenesisAccounts(filename string) []models.GenesisAccount {
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

func getLogConfig() *LogConfig {
	level, err := log.ParseLevel(getString("log.level", "info"))
	if err != nil {
		log.Fatalf("invalid log level: %e", err)
	}

	format := getString("log.format", "text")

	if format != "text" && format != "json" {
		log.Fatalf("invalid log format: %s", format)
	}

	return &LogConfig{
		Level:  level,
		Format: format,
	}
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
