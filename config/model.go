package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/sirupsen/logrus"
)

const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

type Config struct {
	Log       *LogConfig
	Metrics   *MetricsConfig
	Bootstrap *CommanderBootstrapConfig
	Rollup    *RollupConfig
	API       *APIConfig
	Badger    *BadgerConfig
	Ethereum  *EthereumConfig
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

type DeployerConfig struct {
	Bootstrap *DeployerBootstrapConfig
	Ethereum  *EthereumConfig
}

type DeployerBootstrapConfig struct {
	GenesisAccounts  []models.GenesisAccount `json:"-"`
	BlocksToFinalise uint32
}

type CommanderBootstrapConfig struct {
	Prune            bool
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
	RPCURL           string `json:"-"`
	ChainID          uint64
	PrivateKey       string `json:"-"`
	ChainMineTimeout time.Duration
}
