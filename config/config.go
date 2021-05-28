package config

import (
	"log"
	"os"
	"path"
	"strings"

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
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read in config: %s", err)
	}
	return &Config{
		Rollup: &RollupConfig{
			SyncBatches:            viper.GetBool("rollup.sync_batches"),
			FeeReceiverPubKeyID:    viper.GetUint32("rollup.fee_receiver_pub_key_id"),
			TxsPerCommitment:       viper.GetUint32("rollup.txs_per_commitment"),
			MinCommitmentsPerBatch: viper.GetUint32("rollup.min_commitments_per_batch"),
			MaxCommitmentsPerBatch: viper.GetUint32("rollup.max_commitments_per_batch"),
			CommitmentLoopInterval: viper.GetDuration("rollup.commitment_loop_interval"),
			BatchLoopInterval:      viper.GetDuration("rollup.batch_loop_interval"),
			GenesisAccounts:        getGenesisAccounts(),
			BootstrapNodeURL:       getFromViperOrDefault("rollup.bootstrap_node_url", nil),
		},
		API: &APIConfig{
			Version: viper.GetString("api.version"),
			Port:    viper.GetString("api.port"),
			DevMode: viper.GetBool("api.dev_mode"),
		},
		Postgres: &PostgresConfig{
			Host:           getFromViperOrDefault("postgres.host", nil),
			Port:           getFromViperOrDefault("postgres.port", nil),
			Name:           viper.GetString("postgres.name"),
			User:           getFromViperOrDefault("postgres.user", nil),
			Password:       getFromViperOrDefault("postgres.password", nil),
			MigrationsPath: *getFromViperOrDefault("postgres.migrations_path", ref.String(getMigrationsPath())),
		},
		Badger: &BadgerConfig{
			Path: *getFromViperOrDefault("badger.path", ref.String(getBadgerPath())),
		},
		Ethereum: &EthereumConfig{
			RPCURL:     viper.GetString("ethereum.rpc_url"),
			ChainID:    viper.GetString("ethereum.chain_id"),
			PrivateKey: viper.GetString("ethereum.private_key"),
		},
	}
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := *getFromViperOrDefault("rollup.genesis_path", ref.String(getGenesisPath()))

	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Fatalf("error reading genesis file: %s", err.Error())
	}

	return genesisAccounts
}
