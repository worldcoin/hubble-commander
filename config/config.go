package config

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/fsnotify/fsnotify"
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
	viper.SetConfigFile(path.Join(utils.GetProjectRoot(), "config.yml"))
	viper.AutomaticEnv()
	viper.SetEnvPrefix("HUBBLE")

	config := Config{
		Rollup:   &RollupConfig{},
		API:      &APIConfig{},
		Postgres: &PostgresConfig{},
		Badger:   &BadgerConfig{},
		Ethereum: &EthereumConfig{},
	}
	updateConfig(&config)
	config.Rollup.GenesisAccounts = getGenesisAccounts()
	return &config
}

func GetTestConfig() Config {
	return Config{
		Rollup: &RollupConfig{
			SyncBatches:            false,
			FeeReceiverPubKeyID:    0,
			TxsPerCommitment:       2,
			MinCommitmentsPerBatch: 1,
			MaxCommitmentsPerBatch: 32,
			CommitmentLoopInterval: 500 * time.Millisecond,
			BatchLoopInterval:      500 * time.Millisecond,
			GenesisAccounts:        getGenesisAccounts(),
		},
		API: &APIConfig{
			Version: "dev-0.0.1",
			Port:    *getEnvOrDefault("HUBBLE_PORT", ref.String("8080")),
			DevMode: true,
		},
		Postgres: &PostgresConfig{
			Host:           getEnvOrDefault("HUBBLE_DBHOST", nil),
			Port:           getEnvOrDefault("HUBBLE_DBPORT", nil),
			Name:           *getEnvOrDefault("HUBBLE_DBNAME", ref.String("hubble_test")),
			User:           getEnvOrDefault("HUBBLE_DBUSER", nil),
			Password:       getEnvOrDefault("HUBBLE_DBPASSWORD", nil),
			MigrationsPath: getMigrationsPath(),
		},
		Badger: &BadgerConfig{
			Path: *getEnvOrDefault("HUBBLE_BADGER_PATH", ref.String(getBadgerPath())),
		},
		Ethereum: getEthereumConfig(),
	}
}

func updateConfig(cfg *Config) {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read in config: %s", err)
	}

	cfg.Rollup.SyncBatches = viper.GetBool("sync_batches")
	cfg.Rollup.FeeReceiverPubKeyID = viper.GetUint32("fee_receiver_pub_key_id")
	cfg.Rollup.TxsPerCommitment = viper.GetUint32("txs_per_commitment")
	cfg.Rollup.MinCommitmentsPerBatch = viper.GetUint32("min_commitments_per_batch")
	cfg.Rollup.MaxCommitmentsPerBatch = viper.GetUint32("max_commitments_per_batch")
	cfg.Rollup.CommitmentLoopInterval = viper.GetDuration("commitment_loop_interval")
	cfg.Rollup.BatchLoopInterval = viper.GetDuration("batch_loop_interval")

	cfg.API.Version = viper.GetString("version")
	cfg.API.Port = viper.GetString("port")

	cfg.Postgres.Host = getFromViperOrDefault("dbhost", nil)
	cfg.Postgres.Port = getFromViperOrDefault("dbport", nil)
	cfg.Postgres.Name = viper.GetString("dbname")
	cfg.Postgres.User = getFromViperOrDefault("dbuser", nil)
	cfg.Postgres.Password = getFromViperOrDefault("dbpassword", nil)
	cfg.Postgres.MigrationsPath = *getFromViperOrDefault("migrations_path", ref.String(getMigrationsPath()))

	cfg.Badger.Path = *getFromViperOrDefault("badger_path", ref.String(getBadgerPath()))

	viper.SetEnvPrefix("ETHEREUM")
	cfg.Ethereum.RPCURL = viper.GetString("rpc_url")
	cfg.Ethereum.ChainID = viper.GetString("chain_id")
	cfg.Ethereum.PrivateKey = viper.GetString("private_key")
	viper.SetEnvPrefix("HUBBLE")
}

func getGenesisAccounts() []models.GenesisAccount {
	filename := *getEnvOrDefault("HUBBLE_GENESIS_PATH", ref.String(getGenesisPath()))

	genesisAccounts, err := readGenesisFile(filename)
	if err != nil {
		log.Fatalf("error reading genesis file: %s", err.Error())
	}

	return genesisAccounts
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

func WatchConfig(config *Config) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %s", e.Name)
		updateConfig(config)
	})
}
