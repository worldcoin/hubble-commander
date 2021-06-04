package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
)

type Config struct {
	Rollup   *RollupConfig
	API      *APIConfig
	Postgres *PostgresConfig
	Badger   *BadgerConfig
	Ethereum *EthereumConfig
}

type RollupConfig struct {
	// TODO: Extract to a separate BootstrapConfig.
	Prune            bool
	GenesisAccounts  []models.GenesisAccount
	BootstrapNodeURL *string

	SyncBatches bool
	SyncSize    uint32

	FeeReceiverPubKeyID    uint32
	TxsPerCommitment       uint32
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
