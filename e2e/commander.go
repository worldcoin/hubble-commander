package e2e

import (
	"fmt"
	"os"

	"github.com/ybbus/jsonrpc/v2"
)

type Commander interface {
	Start() error
	Stop() error
	Client() jsonrpc.RPCClient
}

func NewCommanderFromEnv() (Commander, error) {
	switch os.Getenv("HUBBLE_E2E") {
	case "", "docker":
		return StartDockerCommander(StartOptions{
			Image: "ghcr.io/worldcoin/hubble-commander:latest",
		})
	case "local":
		return ConnectToLocalCommander(), nil
	default:
		return nil, fmt.Errorf("invalid HUBBLE_E2E env var")
	}
}
