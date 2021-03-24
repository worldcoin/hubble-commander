package config

import "time"

type Config struct {
	Rollup   RollupConfig
	API      APIConfig
	DB       DBConfig
	Ethereum *EthereumConfig
}

type RollupConfig struct {
	FeeReceiverIndex       uint32
	TxsPerCommitment       uint32
	CommitmentLoopInterval time.Duration
}

type APIConfig struct {
	Version string
	Port    string
}

type DBConfig struct {
	Host     *string
	Port     *string
	Name     string
	User     *string
	Password *string
}

type EthereumConfig struct {
	RPCURL     string
	PrivateKey string
	ChainID    string
}
