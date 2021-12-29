package setup

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type Commander interface {
	Start() error
	Stop() error
	Restart() error
	Client() jsonrpc.RPCClient
}

func NewConfiguredCommanderFromEnv(commanderConfig *config.Config, deployerConfig *config.DeployerConfig) (Commander, error) {
	if commanderConfig != nil {
		logRequiredConfig(commanderConfig.Rollup)
	}

	switch os.Getenv("HUBBLE_E2E") {
	case "docker":
		return StartDockerCommander(StartOptions{
			Image:           "ghcr.io/worldcoin/hubble-commander:latest",
			Prune:           true,
			DeployContracts: true,
		})
	case "local":
		return ConnectToLocalCommander(), nil
	case "in-process":
		return DeployAndCreateInProcessCommander(commanderConfig, deployerConfig)
	default:
		return nil, fmt.Errorf("invalid HUBBLE_E2E env var")
	}
}

func logRequiredConfig(cfg *config.RollupConfig) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	log.Printf("Required Rollup config for this test: %s", string(jsonCfg))
}
