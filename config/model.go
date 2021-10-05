package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Log       *LogConfig
	Bootstrap *BootstrapConfig
	Rollup    *RollupConfig
	API       *APIConfig
	Postgres  *PostgresConfig
	Badger    *BadgerConfig
	Ethereum  *EthereumConfig
}

type LogConfig struct {
	Level logrus.Level
	// "json" or "text" (default)
	Format string
}

type BootstrapConfig struct {
	Prune            bool
	GenesisAccounts  []models.GenesisAccount `json:"-"`
	BlocksToFinalise uint32
	BootstrapNodeURL *string
	ChainSpecPath    *string
}

type RollupConfig struct {
	SyncSize                  uint32
	FeeReceiverPubKeyID       uint32
	MinTxsPerCommitment       uint32
	MaxTxsPerCommitment       uint32
	MinCommitmentsPerBatch    uint32
	MaxCommitmentsPerBatch    uint32
	TransitionDisputeGasLimit uint64
	SignatureDisputeGasLimit  uint64
	CommitmentLoopInterval    time.Duration
	BatchLoopInterval         time.Duration
	DisableSignatures         bool
}

type APIConfig struct {
	Version string
	Port    string
}

type PostgresConfig struct {
	Host           *string `json:"-"`
	Port           *string `json:"-"`
	Name           string
	User           *string `json:"-"`
	Password       *string `json:"-"`
	MigrationsPath string
}

type BadgerConfig struct {
	Path string
}

type EthereumConfig struct {
	RPCURL     string `json:"-"`
	PrivateKey string `json:"-"`
	ChainID    string
}
