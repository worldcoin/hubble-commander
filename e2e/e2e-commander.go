package e2e

import (
	"fmt"
	"os"

	"github.com/ybbus/jsonrpc/v2"
)

type E2ECommander interface {
	Start() error
	Stop() error
	Client() jsonrpc.RPCClient
}

func CreateCommanderFromEnv() (E2ECommander, error) {
	switch os.Getenv("HUBBLE_E2E") {
	case "":
		fallthrough
	case "docker":
		return StartCommander(StartOptions{
			Image: "ghcr.io/worldcoin/hubble-commander:latest",
		})
	case "local":
		return ConnectToLocalCommander(), nil
	}

	return nil, fmt.Errorf("invalid HUBBLE_E2E")
}
