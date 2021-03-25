package e2e

import (
	"github.com/ybbus/jsonrpc/v2"
	"os"
	"os/exec"
	"time"
)

type StartOptions struct {
	Image string
}

func StartCommander(opts StartOptions) (*TestCommander, error) {
	cmd := exec.Command(
		"docker", "run", "-t", "--rm",
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

	time.Sleep(1 * time.Second)

	client := jsonrpc.NewClient("http://localhost:8080")

	return &TestCommander{Process: cmd.Process, Client: client}, nil
}

type TestCommander struct {
	Process *os.Process
	Client  jsonrpc.RPCClient
}
