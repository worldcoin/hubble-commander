package config

import (
	"log"
	"os"
	"path"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
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
	return path.Join(utils.GetProjectRoot(), "db", "postgres", "migrations")
}

func getBadgerPath() string {
	return path.Join(utils.GetProjectRoot(), "db", "badger", "data")
}

func getGenesisPath() string {
	return path.Join(utils.GetProjectRoot(), "genesis.yaml")
}

func GetConfig() *Config {
	return newConfig("config.yaml")
}

func GetTestConfig() *Config {
	return newConfig("config-test.yaml")
}

func newConfig(fileName string) *Config {
	viper.SetConfigFile(path.Join(utils.GetProjectRoot(), fileName))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read in config: %s", err)
	}
	cfg := &Config{
		Rollup: &RollupConfig{
			SyncBatches:            viper.GetBool("sync_batches"),
			FeeReceiverPubKeyID:    viper.GetUint32("fee_receiver_pub_key_id"),
			TxsPerCommitment:       viper.GetUint32("txs_per_commitment"),
			MinCommitmentsPerBatch: viper.GetUint32("min_commitments_per_batch"),
			MaxCommitmentsPerBatch: viper.GetUint32("max_commitments_per_batch"),
			CommitmentLoopInterval: viper.GetDuration("commitment_loop_interval"),
			BatchLoopInterval:      viper.GetDuration("batch_loop_interval"),
			GenesisAccounts:        getGenesisAccounts(),
		},
		API: &APIConfig{
			Version: viper.GetString("version"),
			Port:    viper.GetString("port"),
			DevMode: viper.GetBool("dev_mode"),
		},
		Postgres: &PostgresConfig{
			Host:           getFromViperOrDefault("dbhost", nil),
			Port:           getFromViperOrDefault("dbport", nil),
			Name:           viper.GetString("dbname"),
			User:           getFromViperOrDefault("dbuser", nil),
			Password:       getFromViperOrDefault("dbpassword", nil),
			MigrationsPath: *getFromViperOrDefault("migrations_path", ref.String(getMigrationsPath())),
		},
		Badger: &BadgerConfig{
			Path: *getFromViperOrDefault("badger_path", ref.String(getBadgerPath())),
		},
		Ethereum: &EthereumConfig{},
	}

	viper.SetEnvPrefix("ETHEREUM")
	cfg.Ethereum.RPCURL = viper.GetString("rpc_url")
	cfg.Ethereum.ChainID = viper.GetString("chain_id")
	cfg.Ethereum.PrivateKey = viper.GetString("private_key")
	return cfg
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := *getEnvOrDefault("HUBBLE_GENESIS_PATH", ref.String(getGenesisPath()))

	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Fatalf("error reading genesis file: %s", err.Error())
	}

	return genesisAccounts
}
