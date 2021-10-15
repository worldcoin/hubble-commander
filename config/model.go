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
	GenesisAccounts  []models.GenesisAccount `json:"-"`
	BlocksToFinalise uint32
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
	TransitionDisputeGasLimit        uint64
	SignatureDisputeGasLimit         uint64
	BatchAccountRegistrationGasLimit uint64
	BatchLoopInterval                time.Duration
	DisableSignatures                bool
}

type APIConfig struct {
	Version string
	Port    string
}

type BadgerConfig struct {
	Path string
}

type EthereumConfig struct {
	RPCURL     string `json:"-"`
	ChainID    uint64
	PrivateKey string `json:"-"`
}
