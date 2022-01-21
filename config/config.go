package config

import (
	"os"
	"path"
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	SimulatorChainID                        = 1337
	DefaultTransferBatchSubmissionGasLimit  = uint64(400_000)
	DefaultC2TBatchSubmissionGasLimit       = uint64(500_000)
	DefaultMMBatchSubmissionGasLimit        = uint64(550_000)
	DefaultDepositBatchSubmissionGasLimit   = uint64(220_000)
	DefaultTransitionDisputeGasLimit        = uint64(4_100_000)
	DefaultSignatureDisputeGasLimit         = uint64(7_600_000)
	DefaultBatchAccountRegistrationGasLimit = uint64(8_000_000)
	DefaultMetricsPort                      = "2112"
	DefaultMetricsEndpoint                  = "/metrics"
)

func GetConfig() *Config {
	setupViper("commander-config.yaml")

	return &Config{
		Log:     getLogConfig(),
		Metrics: getMetricsConfig(),
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
			TransferBatchSubmissionGasLimit:  getUint64("rollup.transfer_batch_submission_gas_limit", DefaultTransferBatchSubmissionGasLimit),
			C2TBatchSubmissionGasLimit:       getUint64("rollup.c2t_batch_submission_gas_limit", DefaultC2TBatchSubmissionGasLimit),
			MMBatchSubmissionGasLimit:        getUint64("rollup.mm_batch_submission_gas_limit", DefaultMMBatchSubmissionGasLimit),
			DepositBatchSubmissionGasLimit:   getUint64("rollup.deposit_batch_submission_gas_limit", DefaultDepositBatchSubmissionGasLimit),
			TransitionDisputeGasLimit:        getUint64("rollup.transition_dispute_gas_limit", DefaultTransitionDisputeGasLimit),
			SignatureDisputeGasLimit:         getUint64("rollup.signature_dispute_gas_limit", DefaultSignatureDisputeGasLimit),
			BatchAccountRegistrationGasLimit: getUint64("rollup.batch_account_registration_gas_limit", DefaultBatchAccountRegistrationGasLimit),
			BatchLoopInterval:                getDuration("rollup.batch_loop_interval", 500*time.Millisecond),
			DisableSignatures:                getBool("rollup.disable_signatures", false),
			MaxTxnDelay:                      getDuration("rollup.max_txn_delay", 30*time.Minute),
		},
		API: &APIConfig{
			Version:            "0.5.0-rc2",
			Port:               getString("api.port", "8080"),
			EnableProofMethods: getBool("api.enable_proof_methods", false),
			AuthenticationKey:  getStringOrNil("api.enable_proof_methods"),
		},
		Badger: &BadgerConfig{
			Path: getString("badger.path", "./db/data/hubble"),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetTestConfig() *Config {
	setupViper("commander-config.yaml")

	return &Config{
		Log: &LogConfig{
			Level:  log.InfoLevel,
			Format: LogFormatText,
		},
		Metrics: &MetricsConfig{
			Port:     DefaultMetricsPort,
			Endpoint: DefaultMetricsEndpoint,
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
			TransferBatchSubmissionGasLimit:  DefaultTransferBatchSubmissionGasLimit,
			C2TBatchSubmissionGasLimit:       DefaultC2TBatchSubmissionGasLimit,
			MMBatchSubmissionGasLimit:        DefaultMMBatchSubmissionGasLimit,
			DepositBatchSubmissionGasLimit:   DefaultDepositBatchSubmissionGasLimit,
			TransitionDisputeGasLimit:        DefaultTransitionDisputeGasLimit,
			SignatureDisputeGasLimit:         DefaultSignatureDisputeGasLimit,
			BatchAccountRegistrationGasLimit: DefaultBatchAccountRegistrationGasLimit,
			BatchLoopInterval:                500 * time.Millisecond,
			DisableSignatures:                true,
			MaxTxnDelay:                      30 * time.Minute,
		},
		API: &APIConfig{
			Version:            "dev-0.5.0-rc2",
			Port:               "8080",
			EnableProofMethods: true,
			AuthenticationKey:  ref.String("secret_authentication_key"),
		},
		Badger: &BadgerConfig{
			Path: "./db/data/hubble_test",
		},
		Ethereum: &EthereumConfig{
			RPCURL:     "simulator",
			ChainID:    SimulatorChainID,
			PrivateKey: "ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82",
		},
	}
}

func setupViper(configName string) {
	// Find the config file
	viper.SetConfigName(configName)
	viper.AddConfigPath("/etc/hubble")
	viper.AddConfigPath("$HOME/.hubble")
	viper.AddConfigPath(".") // Current working dir

	// Add the *parent* of folder containing the executable
	// So when it's run from `some_dir/build/hubble` we will look in `some_dir`
	// TODO: Remove
	exe_path, err := os.Executable()
	if err != nil {
		log.Panic("Could not find exe location: %s", err)
	}
	viper.AddConfigPath(path.Join(path.Dir(exe_path), ".."))

	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err = viper.ReadInConfig()
	if err != nil {
		if strings.Contains(err.Error(), "Not Found in") {
			log.Warn(err)
			log.Warn("Continuing with default config (possibly overridden by env vars).")
		} else {
			log.Panicf("failed to read in config: %s", err)
		}
	}
}

func getLogConfig() *LogConfig {
	level, err := log.ParseLevel(getString("log.level", "info"))
	if err != nil {
		log.Panicf("invalid log level: %e", err)
	}

	format := getString("log.format", LogFormatText)

	if format != LogFormatText && format != LogFormatJSON {
		log.Panicf("invalid log format: %s", format)
	}

	return &LogConfig{
		Level:  level,
		Format: format,
	}
}

func getMetricsConfig() *MetricsConfig {
	return &MetricsConfig{
		Port:     getString("metrics.port", DefaultMetricsPort),
		Endpoint: getString("metrics.endpoint", DefaultMetricsEndpoint),
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
		ChainID:    getUint64OrPanic("ethereum.chain_id"),
		PrivateKey: getStringOrPanic("ethereum.private_key"),
	}
}
