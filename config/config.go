package config

import (
	"os"
	"path"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/joho/godotenv"
)

var (
	genesisAccounts = []models.GenesisAccount{
		{
			PrivateKey: [32]byte{47, 122, 85, 155, 45, 45, 78, 193, 227, 186, 188, 1, 34, 231, 239, 12,
				106, 69, 205, 180, 204, 209, 103, 244, 86, 202, 202, 82, 17, 35, 254, 158},
			Balance: models.MakeUint256(1000),
		},
		{
			PrivateKey: [32]byte{1, 49, 177, 240, 47, 37, 4, 166, 10, 48, 38, 31, 163, 102, 92, 161, 46,
				213, 66, 171, 1, 247, 61, 235, 177, 155, 175, 169, 150, 121, 2, 114},
			Balance: models.MakeUint256(1000),
		},
		{
			PrivateKey: [32]byte{10, 9, 162, 211, 112, 191, 164, 2, 77, 121, 49, 230, 55, 131, 232, 78,
				138, 60, 51, 46, 182, 19, 63, 38, 141, 234, 171, 114, 217, 112, 209, 102},
			Balance: models.MakeUint256(1000),
		},
		{
			PrivateKey: [32]byte{1, 49, 177, 240, 47, 37, 4, 166, 10, 48, 38, 31, 163, 102, 92, 161, 46,
				213, 66, 171, 1, 247, 61, 235, 177, 155, 175, 169, 150, 121, 2, 114},
			Balance: models.MakeUint256(1000),
		},
		{
			PrivateKey: [32]byte{28, 144, 144, 86, 206, 85, 56, 204, 99, 151, 175, 222, 43, 236, 189, 2, 69, 132, 135, 164, 229, 121, 134, 181, 231, 109, 56, 204, 16, 81, 134, 58},
		},
		{
			PrivateKey: [32]byte{33, 167, 159, 167, 55, 75, 126, 104, 141, 124, 210, 92, 208, 195, 87, 114, 79, 18, 225, 124, 61, 170, 42, 128, 231, 29, 48, 12, 37, 208, 219, 120},
		},
	}
)

func init() {
	if os.Getenv("CI") != "true" {
		loadDotEnv()
	}
}

func loadDotEnv() {
	dotEnvFilePath := path.Join(utils.GetProjectRoot(), ".env")
	_ = godotenv.Load(dotEnvFilePath)
}

func getMigrationsPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "migrations")
}

func GetConfig() Config {
	return Config{
		Rollup: RollupConfig{
			FeeReceiverIndex:        0,
			TxsPerCommitment:        2,
			MinCommitmentsPerBatch:  1,
			MaxCommitmentsPerBatch:  32,
			CommitmentLoopInterval:  500 * time.Millisecond,
			BatchLoopInterval:       500 * time.Millisecond,
			BlockNumberLoopInterval: 500 * time.Millisecond,
			GenesisAccounts:         genesisAccounts,
			SignaturesDomain:        [32]byte{1, 2, 3, 4},
		},
		API: APIConfig{
			Version: "dev-0.1.0",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: false,
		},
		DB: DBConfig{
			Host:           getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:           getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:           *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble")),
			User:           getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password:       getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
			MigrationsPath: *getEnvOrDefault("HUBBLE_MIGRATIONS_PATH", ref.String(getMigrationsPath())),
		},
		Ethereum: getEthereumConfig(),
	}
}

func GetTestConfig() Config {
	return Config{
		Rollup: RollupConfig{
			FeeReceiverIndex:        0,
			TxsPerCommitment:        2,
			MinCommitmentsPerBatch:  1,
			MaxCommitmentsPerBatch:  32,
			CommitmentLoopInterval:  500 * time.Millisecond,
			BatchLoopInterval:       500 * time.Millisecond,
			BlockNumberLoopInterval: 500 * time.Millisecond,
			GenesisAccounts:         genesisAccounts,
			SignaturesDomain:        [32]byte{1, 2, 3, 4},
		},
		API: APIConfig{
			Version: "dev-0.1.0",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: false,
		},
		DB: DBConfig{
			Host:           getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:           getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:           *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble_test")),
			User:           getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password:       getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
			MigrationsPath: getMigrationsPath(),
		},
		Ethereum: getEthereumConfig(),
	}
}

func getEthereumConfig() *EthereumConfig {
	rpcURL := getEnvOrDefault("ETHEREUM_RPC_URL", nil)
	if rpcURL == nil {
		return nil
	}
	return &EthereumConfig{
		RPCURL:     *rpcURL,
		ChainID:    getEnv("ETHEREUM_CHAIN_ID"),
		PrivateKey: getEnv("ETHEREUM_PRIVATE_KEY"),
	}
}
