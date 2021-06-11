package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
)

type Config struct {
	Bootstrap *BootstrapConfig
	Rollup    *RollupConfig
	API       *APIConfig
	Postgres  *PostgresConfig
	Badger    *BadgerConfig
	Ethereum  *EthereumConfig
}

type BootstrapConfig struct {
	Prune            bool
	GenesisAccounts  []models.GenesisAccount
	BootstrapNodeURL *string
}

type RollupConfig struct {
	SyncSize               uint32
	FeeReceiverPubKeyID    uint32
	TxsPerCommitment       uint64
	MinCommitmentsPerBatch uint32
	MaxCommitmentsPerBatch uint32
	CommitmentLoopInterval time.Duration
	BatchLoopInterval      time.Duration
}

type APIConfig struct {
	Version string
	Port    string
	DevMode bool
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
