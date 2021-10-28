package config

import (
	"path"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	SimulatorChainID                        = 1337
	DefaultTransitionDisputeGasLimit        = uint64(5_000_000)
	DefaultSignatureDisputeGasLimit         = uint64(7_500_000)
	DefaultBatchAccountRegistrationGasLimit = uint64(8_000_000)
	DefaultBlocksToFinalise                 = uint32(7 * 24 * 60 * 4)
)

func setupViperForCommander() {
	setupViper(getCommanderConfigPath())
}

func setupViperForDeployer() {
	setupViper(getDeployerConfigPath())
}

func setupViper(configPath string) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Printf("Configuration file not found (%s). Continuing with default config (possibly overridden by env vars).", configPath)
		} else {
			log.Fatalf("failed to read in config: %s", err)
		}
	}
}

func GetConfig() *Config {
	setupViperForCommander()

	return &Config{
		Log: getLogConfig(),
		Bootstrap: &CommanderBootstrapConfig{
			Prune:            getBool("bootstrap.prune", false),
			BootstrapNodeURL: getStringOrNil("bootstrap.node_url"),
			ChainSpecPath:    getStringOrNil("bootstrap.chain_spec_path"),
		},
		Rollup: &RollupConfig{
			SyncSize:                         getUint32("rollup.sync_size", 50),
			FeeReceiverPubKeyID:              getUint32("rollup.fee_receiver_pub_key_id", 0),
			MinTxsPerCommitment:              getUint32("rollup.min_txs_per_commitment", 1),
			MaxTxsPerCommitment:              getUint32("rollup.max_txs_per_commitment", 32),
			MinCommitmentsPerBatch:           getUint32("rollup.min_commitments_per_batch", 1),
			MaxCommitmentsPerBatch:           getUint32("rollup.max_commitments_per_batch", 32),
			TransitionDisputeGasLimit:        getUint64("rollup.transition_dispute_gas_limit", DefaultTransitionDisputeGasLimit),
			SignatureDisputeGasLimit:         getUint64("rollup.signature_dispute_gas_limit", DefaultSignatureDisputeGasLimit),
			BatchAccountRegistrationGasLimit: getUint64("rollup.batch_account_registration_gas_limit", DefaultBatchAccountRegistrationGasLimit),
			BatchLoopInterval:                getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			DisableSignatures:                getBool("rollup.disable_signatures", false),
		},
		API: &APIConfig{
			Version:            "0.5.0-rc2",
			Port:               getString("api.port", "8080"),
			EnableProofMethods: getBool("api.enable_proof_methods", false),
		},
		Badger: &BadgerConfig{
			Path: getString("badger.path", getBadgerPath()),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetTestConfig() *Config {
	setupViperForCommander()

	return &Config{
		Log: &LogConfig{
			Level:  log.InfoLevel,
			Format: LogFormatText,
		},
		Bootstrap: &CommanderBootstrapConfig{
			Prune:            false,
			BootstrapNodeURL: nil,
			ChainSpecPath:    nil,
		},
		Rollup: &RollupConfig{
			SyncSize:                         50,
			FeeReceiverPubKeyID:              0,
			MinTxsPerCommitment:              2,
			MaxTxsPerCommitment:              2,
			MinCommitmentsPerBatch:           1,
			MaxCommitmentsPerBatch:           32,
			TransitionDisputeGasLimit:        DefaultTransitionDisputeGasLimit,
			SignatureDisputeGasLimit:         DefaultSignatureDisputeGasLimit,
			BatchAccountRegistrationGasLimit: DefaultBatchAccountRegistrationGasLimit,
			BatchLoopInterval:                500 * time.Millisecond,
			DisableSignatures:                true,
		},
		API: &APIConfig{
			Version: "dev-0.5.0-rc2",
			Port:    "8080",
		},
		Badger: &BadgerConfig{
			Path: getTestBadgerPath(),
		},
		Ethereum: &EthereumConfig{
			RPCURL:     "simulator",
			ChainID:    SimulatorChainID,
			PrivateKey: "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82",
		},
	}
}

func GetDeployerConfig() *DeployerConfig {
	setupViperForDeployer()

	return &DeployerConfig{
		Bootstrap: &DeployerBootstrapConfig{
			BlocksToFinalise: getUint32("bootstrap.blocks_to_finalise", DefaultBlocksToFinalise), // nolint:misspell
			GenesisAccounts:  getGenesisAccounts(),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getCommanderConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "commander-config.yaml")
}

func getDeployerConfigPath() string {
	return path.Join(utils.GetProjectRoot(), "deployer-config.yaml")
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

func getBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "data", "hubble")
}

func getTestBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "data", "hubble_test")
}

func getLogConfig() *LogConfig {
	level, err := log.ParseLevel(getString("log.level", "info"))
	if err != nil {
		log.Fatalf("invalid log level: %e", err)
	}

	format := getString("log.format", LogFormatText)

	if format != LogFormatText && format != LogFormatJSON {
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
			ChainID:    SimulatorChainID,
			PrivateKey: getString("ethereum.private_key", "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82"),
		}
	}
	return &EthereumConfig{
		RPCURL:     *rpcURL,
		ChainID:    getUint64OrThrow("ethereum.chain_id"),
		PrivateKey: getStringOrThrow("ethereum.private_key"),
	}
}
