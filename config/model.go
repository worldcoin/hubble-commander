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
	BootstrapNodeURL *string                 // TODO-CHAIN - remove
	ChainSpecPath    *string
}

type RollupConfig struct {
	SyncSize               uint32
	FeeReceiverPubKeyID    uint32
	MinTxsPerCommitment    uint32
	MaxTxsPerCommitment    uint32
	MinCommitmentsPerBatch uint32
	MaxCommitmentsPerBatch uint32
	CommitmentLoopInterval time.Duration
	BatchLoopInterval      time.Duration
	DisableSignatures      bool
}

type APIConfig struct {
	Version string
	Port    string
}

type PostgresConfig struct {
	Host           *string
	Port           *string
	Name           string
	User           *string
	Password       *string
	MigrationsPath string
}

type BadgerConfig struct {
	Path string
}

type EthereumConfig struct {
	RPCURL     string
	PrivateKey string
	ChainID    string
}
