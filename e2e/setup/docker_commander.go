package setup

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Worldcoin/hubble-commander/commander"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"
	"github.com/ybbus/jsonrpc/v2"
)

type DockerCommander struct {
	cli         *docker.Client
	containerID string
	client      jsonrpc.RPCClient
}

type StartOptions struct {
	Image           string
	Prune           bool
	DeployContracts bool
}

func getEnvVarMapping(key string) string {
	return key + "=" + os.Getenv(key)
}

func StartDockerCommander(opts StartOptions) (*DockerCommander, error) {
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return nil, err
	}

	var networkMode container.NetworkMode
	if os.Getenv("CI") == "true" {
		networkMode = "host"
	} else {
		networkMode = "bridge"
	}

	if opts.DeployContracts {
		err = deployContractsAndStoreChainSpec()
		if err != nil {
			return nil, err
		}
	}

	created, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: opts.Image,
			Env: []string{
				getEnvVarMapping("HUBBLE_ETHEREUM_RPC_URL"),
				getEnvVarMapping("HUBBLE_ETHEREUM_CHAIN_ID"),
				getEnvVarMapping("HUBBLE_ETHEREUM_PRIVATE_KEYS"),
				getEnvVarMapping("HUBBLE_ROLLUP_MIN_TXS_PER_COMMITMENT"),
				getEnvVarMapping("HUBBLE_ROLLUP_MAX_TXS_PER_COMMITMENT"),
				"HUBBLE_API_PORT=8080",
				"HUBBLE_BOOTSTRAP_CHAIN_SPEC_PATH=chain-spec/chain-spec.yaml",
				fmt.Sprintf("HUBBLE_BOOTSTRAP_PRUNE=%t", opts.Prune),
			},
			ExposedPorts: map[nat.Port]struct{}{
				"8080/tcp": {},
			},
		},
		&container.HostConfig{
			NetworkMode: networkMode,
			PortBindings: map[nat.Port][]nat.PortBinding{
				"8080/tcp": {
					nat.PortBinding{HostIP: "", HostPort: "8080"},
				},
			},
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: path.Join(".", "e2e-data"),
					Target: "/go/src/app/db/data/hubble",
				},
				{
					Type:   mount.TypeBind,
					Source: path.Join(".", "e2e-chain-spec"),
					Target: "/go/src/app/chain-spec",
				},
			},
		},
		&network.NetworkingConfig{},
		"",
	)
	if err != nil {
		return nil, err
	}

	containerID := created.ID

	client := jsonrpc.NewClient("http://localhost:8080")

	cmd := &DockerCommander{
		cli:         cli,
		containerID: containerID,
		client:      client,
	}

	stream, err := cli.ContainerAttach(context.Background(), containerID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true,
	})
	if err != nil {
		return nil, err
	}
	go func() {
		_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, stream.Reader)
		if err != nil {
			log.Panicf("%+v", err)
		}
	}()

	return cmd, nil
}

func deployContractsAndStoreChainSpec() error {
	deployerCfg := config.GetDeployerTestConfig()

	blockchain, err := commander.GetChainConnection(deployerCfg.Ethereum)
	if err != nil {
		return err
	}

	chainSpec, err := commander.Deploy(deployerCfg, blockchain)
	if err != nil {
		return err
	}

	chainSpecPath := path.Join(".", "e2e-chain-spec", "chain-spec.yaml")
	return utils.StoreChainSpec(chainSpecPath, *chainSpec)
}

func (c *DockerCommander) Client() jsonrpc.RPCClient {
	return c.client
}

func (c *DockerCommander) Start() error {
	err := c.cli.ContainerStart(context.Background(), c.containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	start := time.Now()
	for {
		healthy, err := c.IsHealthy()
		if err != nil {
			return err
		}

		if healthy {
			break
		}

		hasExited, err := c.HasExited()
		if err != nil {
			return err
		}

		if hasExited {
			return fmt.Errorf("container has exited")
		}

		if time.Since(start) > 120*time.Second {
			return fmt.Errorf("node start timeout")
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (c *DockerCommander) IsHealthy() (bool, error) {
	info, err := c.cli.ContainerInspect(context.Background(), c.containerID)
	if err != nil {
		return false, err
	}

	return info.State != nil && info.State.Health != nil && info.State.Health.Status == "healthy", nil
}

func (c *DockerCommander) HasExited() (bool, error) {
	info, err := c.cli.ContainerInspect(context.Background(), c.containerID)
	if err != nil {
		return false, err
	}

	return info.State != nil && info.State.Status == "exited", nil
}

func (c *DockerCommander) Stop() error {
	err := c.cli.ContainerStop(context.Background(), c.containerID, nil)
	if err != nil {
		return err
	}
	err = c.cli.ContainerRemove(context.Background(), c.containerID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		return err
	}

	return nil
}

func (c *DockerCommander) Restart() error {
	err := c.Stop()
	if err != nil {
		return err
	}

	cmd, err := StartDockerCommander(StartOptions{
		Image:           "ghcr.io/worldcoin/hubble-commander:latest",
		Prune:           false,
		DeployContracts: false,
	})
	if err != nil {
		return err
	}

	c.containerID = cmd.containerID
	return cmd.Start()
}
