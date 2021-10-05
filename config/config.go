package config

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	SimulatorChainID                 = 1337
	DefaultTransitionDisputeGasLimit = uint64(5_000_000)
	DefaultSignatureDisputeGasLimit  = uint64(7_500_000)
	DefaultBlocksToFinalise          = uint32(7 * 24 * 60 * 4)
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
			Prune:            getBool("bootstrap.prune", false),
			GenesisAccounts:  getGenesisAccounts(),
			BlocksToFinalise: getUint32("bootstrap.blocks_to_finalise", DefaultBlocksToFinalise), // nolint:misspell
			BootstrapNodeURL: getStringOrNil("bootstrap.node_url"),
			ChainSpecPath:    getStringOrNil("bootstrap.chain_spec_path"),
		},
		Rollup: &RollupConfig{
			SyncSize:                  getUint32("rollup.sync_size", 50),
			FeeReceiverPubKeyID:       getUint32("rollup.fee_receiver_pub_key_id", 0),
			MinTxsPerCommitment:       getUint32("rollup.min_txs_per_commitment", 1),
			MaxTxsPerCommitment:       getUint32("rollup.max_txs_per_commitment", 32),
			MinCommitmentsPerBatch:    getUint32("rollup.min_commitments_per_batch", 1),
			MaxCommitmentsPerBatch:    getUint32("rollup.max_commitments_per_batch", 32),
			TransitionDisputeGasLimit: getUint64("rollup.transition_dispute_gas_limit", DefaultTransitionDisputeGasLimit),
			SignatureDisputeGasLimit:  getUint64("rollup.signature_dispute_gas_limit", DefaultSignatureDisputeGasLimit),
			CommitmentLoopInterval:    getDuration("rollup.commitment_loop_interval", 500*time.Millisecond),
			BatchLoopInterval:         getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			DisableSignatures:         getBool("rollup.disable_signatures", false),
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
	setupViper()

	return &Config{
		Log: &LogConfig{
			Level:  log.InfoLevel,
			Format: "text",
		},
		Bootstrap: &BootstrapConfig{
			Prune:            false,
			GenesisAccounts:  readGenesisAccounts(getGenesisPath()),
			BlocksToFinalise: DefaultBlocksToFinalise,
			BootstrapNodeURL: nil,
			ChainSpecPath:    nil,
		},
		Rollup: &RollupConfig{
			SyncSize:                  50,
			FeeReceiverPubKeyID:       0,
			MinTxsPerCommitment:       2,
			MaxTxsPerCommitment:       2,
			MinCommitmentsPerBatch:    1,
			MaxCommitmentsPerBatch:    32,
			TransitionDisputeGasLimit: DefaultTransitionDisputeGasLimit,
			SignatureDisputeGasLimit:  DefaultSignatureDisputeGasLimit,
			CommitmentLoopInterval:    500 * time.Millisecond,
			BatchLoopInterval:         500 * time.Millisecond,
			DisableSignatures:         true,
		},
		API: &APIConfig{
			Version: "dev-0.0.1",
			Port:    "8080",
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
			Path: getTestBadgerPath(),
		},
		Ethereum: &EthereumConfig{
			RPCURL:     "simulator",
			ChainID:    strconv.Itoa(SimulatorChainID),
			PrivateKey: "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82",
		},
	}
}

func getConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "config.yaml")
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := getString("bootstrap.genesis_path", getGenesisPath())
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
	return path.Join(utils.GetProjectRoot(), "db", "badger", "data", "hubble")
}

func getTestBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "badger", "data", "hubble_test")
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
		return &EthereumConfig{
			RPCURL:     "simulator",
			ChainID:    strconv.Itoa(SimulatorChainID),
			PrivateKey: getString("ethereum.private_key", "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82"),
		}
	}
	return &EthereumConfig{
		RPCURL:     *rpcURL,
		ChainID:    getStringOrThrow("ethereum.chain_id"),
		PrivateKey: getStringOrThrow("ethereum.private_key"),
	}
}
