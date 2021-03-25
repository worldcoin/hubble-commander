package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/ybbus/jsonrpc/v2"
)

type StartOptions struct {
	Image             string
	UseHostNetworking bool
}

func StartCommander(opts StartOptions) (*TestCommander, error) {
	var networkType string
	if opts.UseHostNetworking {
		networkType = "host"
	} else {
		networkType = "bridge"
	}

	cmd := exec.Command(
		"docker", "run", "-t", "--rm",
		"--network", networkType,
		"-p", "8080:8080",
		"-e", "ETHEREUM_RPC_URL",
		"-e", "ETHEREUM_CHAIN_ID",
		"-e", "ETHEREUM_PRIVATE_KEY",
		"-e", "HUBBLE_DBHOST",
		"-e", "HUBBLE_DBPORT",
		"-e", "HUBBLE_DBNAME",
		"-e", "HUBBLE_DBUSER",
		"-e", "HUBBLE_DBPASSWORD",
		"-e", "HUBBLE_PORT",
		opts.Image,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	client := jsonrpc.NewClient("http://localhost:8080")

	for {
		var version string
		err = client.CallFor(&version, "hubble_getVersion", []interface{}{})
		if err == nil {
			return &TestCommander{Process: cmd.Process, Client: client}, nil
		}
		fmt.Printf("%s\n", err.Error())
		time.Sleep(1 * time.Second)
	}
}

type TestCommander struct {
	Process *os.Process
	Client  jsonrpc.RPCClient
}
