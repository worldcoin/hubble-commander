package config

import (
	"time"

	"github.com/Worldcoin/hubble-commander/models"
)

type Config struct {
	Rollup   RollupConfig
	API      APIConfig
	DB       DBConfig
	Ethereum *EthereumConfig
}

type RollupConfig struct {
	FeeReceiverIndex        uint32
	TxsPerCommitment        uint32
	MinCommitmentsPerBatch  uint32
	MaxCommitmentsPerBatch  uint32
	CommitmentLoopInterval  time.Duration
	BatchLoopInterval       time.Duration
	BlockNumberLoopInterval time.Duration
	GenesisAccounts         []models.GenesisAccount
	SignaturesDomain        [32]byte
}

type APIConfig struct {
	Version string
	Port    string
	DevMode bool
}

type DBConfig struct {
	Host           *string
	Port           *string
	Name           string
	User           *string
	Password       *string
	MigrationsPath string
}

type EthereumConfig struct {
	RPCURL     string
	PrivateKey string
	ChainID    string
}
