package e2e

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/ybbus/jsonrpc/v2"
)

type StartOptions struct {
	Image string
}

func StartCommander(opts StartOptions) (*TestCommander, error) {
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

	created, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: opts.Image,
			Env: []string{
				"ETHEREUM_RPC_URL=" + os.Getenv("ETHEREUM_RPC_URL"),
				"ETHEREUM_CHAIN_ID=" + os.Getenv("ETHEREUM_CHAIN_ID"),
				"ETHEREUM_PRIVATE_KEY=" + os.Getenv("ETHEREUM_PRIVATE_KEY"),
				"HUBBLE_DBHOST=" + os.Getenv("HUBBLE_DBHOST"),
				"HUBBLE_DBPORT=" + os.Getenv("HUBBLE_DBPORT"),
				"HUBBLE_DBNAME=" + os.Getenv("HUBBLE_DBNAME"),
				"HUBBLE_DBUSER=" + os.Getenv("HUBBLE_DBUSER"),
				"HUBBLE_DBPASSWORD=" + os.Getenv("HUBBLE_DBPASSWORD"),
				"HUBBLE_PORT=8080",
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
		},
		&network.NetworkingConfig{},
		"",
	)
	if err != nil {
		return nil, err
	}

	containerID := created.ID

	client := jsonrpc.NewClient("http://localhost:8080")

	commander := &TestCommander{
		cli:         cli,
		containerID: containerID,
		Client:      client,
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
			log.Fatal(err)
		}
	}()

	err = commander.Start()
	if err != nil {
		return nil, err
	}

	return commander, nil
}

type TestCommander struct {
	cli         *docker.Client
	containerID string
	Client      jsonrpc.RPCClient
}

func (c *TestCommander) Start() error {
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

		if time.Since(start) > 30*time.Second {
			return fmt.Errorf("node start timeout")
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func (c *TestCommander) IsHealthy() (bool, error) {
	info, err := c.cli.ContainerInspect(context.Background(), c.containerID)
	if err != nil {
		return false, err
	}

	return info.State != nil && info.State.Health != nil && info.State.Health.Status == "healthy", nil
}

func (c *TestCommander) HasExited() (bool, error) {
	info, err := c.cli.ContainerInspect(context.Background(), c.containerID)
	if err != nil {
		return false, err
	}

	return info.State != nil && info.State.Status == "exited", nil
}

func (c *TestCommander) Stop() error {
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
