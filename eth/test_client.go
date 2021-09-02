package eth

import (
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
)

type TestClient struct {
	*Client
	*simulator.Simulator
}

// NewTestClient Sets up a TestClient backed by automining simulator.
// Remember to call Close() at the end of the test
func NewTestClient() (*TestClient, error) {
	return NewConfiguredTestClient(rollup.DeploymentConfig{}, ClientConfig{})
}

func NewConfiguredTestClient(cfg rollup.DeploymentConfig, clientCfg ClientConfig) (*TestClient, error) {
	sim, err := simulator.NewAutominingSimulator()
	if err != nil {
		return nil, err
	}

	contracts, err := rollup.DeployConfiguredRollup(sim, cfg)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(sim, &NewClientParams{
		ChainState: models.ChainState{
			Rollup: contracts.RollupAddress,
		},
		Rollup:             contracts.Rollup,
		AccountRegistry:    contracts.AccountRegistry,
		TokenRegistry:      contracts.TokenRegistry,
		CustomTokenAddress: contracts.CustomTokenAddress,
		ClientConfig:       clientCfg,
	})
	if err != nil {
		return nil, err
	}

	return &TestClient{
		Client:    client,
		Simulator: sim,
	}, nil
}
