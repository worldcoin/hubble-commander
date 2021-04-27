// nolint:lll
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
	oneEther        = models.MakeUint256FromBig(*utils.ParseEther("1"))
	genesisAccounts = []models.GenesisAccount{
		{
			// PublicKey: 0x0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88
			PrivateKey: [32]byte{47, 122, 85, 155, 45, 45, 78, 193, 227, 186, 188, 1, 34, 231, 239, 12, 106, 69, 205, 180, 204, 209, 103, 244, 86, 202, 202, 82, 17, 35, 254, 158},
			Balance:    oneEther,
		},
		{
			// PublicKey: 0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc
			PrivateKey: [32]byte{1, 49, 177, 240, 47, 37, 4, 166, 10, 48, 38, 31, 163, 102, 92, 161, 46, 213, 66, 171, 1, 247, 61, 235, 177, 155, 175, 169, 150, 121, 2, 114},
			Balance:    oneEther,
		},
		{
			// PublicKey: 0x1ccf19871320b7e850475845d879a9f9717a6c9694fab19498e4261b442de4e011406bdc967984771508a2e50d774f49db36bf5b04b15f9f411b8c8733fe0d8e301f8f2e9aa98f7dde7de3635baa216fdc969e752f4ef646fd5f81d89e46d39804c0ac92c7ea4cc5957b4214ef41a0aa4f1a6f343cebfb577e9dcaf8ff2551d5
			PrivateKey: [32]byte{10, 9, 162, 211, 112, 191, 164, 2, 77, 121, 49, 230, 55, 131, 232, 78, 138, 60, 51, 46, 182, 19, 63, 38, 141, 234, 171, 114, 217, 112, 209, 102},
			Balance:    oneEther,
		},
		{
			// PublicKey: 0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc
			PrivateKey: [32]byte{1, 49, 177, 240, 47, 37, 4, 166, 10, 48, 38, 31, 163, 102, 92, 161, 46, 213, 66, 171, 1, 247, 61, 235, 177, 155, 175, 169, 150, 121, 2, 114},
			Balance:    oneEther,
		},
		{
			PrivateKey: [32]byte{28, 144, 144, 86, 206, 85, 56, 204, 99, 151, 175, 222, 43, 236, 189, 2,
				69, 132, 135, 164, 229, 121, 134, 181, 231, 109, 56, 204, 16, 81, 134, 58},
		},
		{
			PrivateKey: [32]byte{33, 167, 159, 167, 55, 75, 126, 104, 141, 124, 210, 92, 208, 195, 87, 114,
				79, 18, 225, 124, 61, 170, 42, 128, 231, 29, 48, 12, 37, 208, 219, 120},
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
			TxsPerCommitment:        32,
			MinCommitmentsPerBatch:  1,
			MaxCommitmentsPerBatch:  2,
			CommitmentLoopInterval:  500 * time.Millisecond,
			BatchLoopInterval:       500 * time.Millisecond,
			BlockNumberLoopInterval: 500 * time.Millisecond,
			GenesisAccounts:         genesisAccounts,
			SignaturesDomain:        [32]byte{1, 2, 3, 4},
		},
		API: APIConfig{
			Version: "0.0.1",
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
			Version: "dev-0.0.1",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: true,
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
