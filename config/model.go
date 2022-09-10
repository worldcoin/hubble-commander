package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

type Config struct {
	Log       *LogConfig
	Metrics   *MetricsConfig
	Tracing   *TracingConfig
	Bootstrap *CommanderBootstrapConfig
	Rollup    *RollupConfig
	API       *APIConfig
	Badger    *BadgerConfig
	Ethereum  *EthereumConfig

	// Hubble is not yet stable but a lot of services rely on the commander being available
	// at all times. When SafeMode=true Hubble only serves API requests, it does not attempt
	// to create batches or sync against the chain.
	// export HUBBLE_SAFE_MODE=true
	SafeMode bool
}

type LogConfig struct {
	Level logrus.Level
	// "json" or "text" (default)
	Format string
}

type MetricsConfig struct {
	Port     string
	Endpoint string
}

type TracingConfig struct {
	Endpoint    string
	ServiceName string
	Version     string
	Env         string
	Enabled     bool
}

type DeployerConfig struct {
	Bootstrap *DeployerBootstrapConfig
	Ethereum  *EthereumConfig
}

type DeployerBootstrapConfig struct {
	GenesisAccounts  []models.GenesisAccount `json:"-"`
	BlocksToFinalise uint32
	Chooser          *common.Address
}

type CommanderBootstrapConfig struct {
	Prune            bool
	Migrate          bool
	BootstrapNodeURL *string
	ChainSpecPath    *string
}

type RollupConfig struct {
	SyncSize                         uint32
	FeeReceiverPubKeyID              uint32
	MinTxsPerCommitment              uint32
	MaxTxsPerCommitment              uint32
	MinCommitmentsPerBatch           uint32
	MaxCommitmentsPerBatch           uint32
	TransferBatchSubmissionGasLimit  uint64
	C2TBatchSubmissionGasLimit       uint64
	MMBatchSubmissionGasLimit        uint64
	DepositBatchSubmissionGasLimit   uint64
	TransitionDisputeGasLimit        uint64
	SignatureDisputeGasLimit         uint64
	BatchAccountRegistrationGasLimit uint64
	StakeWithdrawalGasLimit          uint64
	BatchLoopInterval                time.Duration
	DisableSignatures                bool

	// if MinTxsPerCommitment or MinCommitmentsPerBatch cause a pending transaction to
	// wait to be included for longer than this delay then they will be ignored and a
	// new batch will be submitted
	MaxTxnDelay time.Duration
}

type APIConfig struct {
	Version            string
	Port               string
	EnableProofMethods bool
	AuthenticationKey  string `json:"-"`
}

type BadgerConfig struct {
	Path string
}

type EthereumConfig struct {
	RPCURL      string `json:"-"`
	ChainID     uint64
	PrivateKey  string `json:"-"`
	MineTimeout time.Duration
}
