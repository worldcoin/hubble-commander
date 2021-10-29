package setup

import (
	"fmt"
	"os"

	"github.com/ybbus/jsonrpc/v2"
)

type Commander interface {
	Start() error
	Stop() error
	Restart() error
	Client() jsonrpc.RPCClient
}

func NewCommanderFromEnv() (Commander, error) {
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
		return CreateInProcessCommander()
	default:
		return nil, fmt.Errorf("invalid HUBBLE_E2E env var")
	}
}
