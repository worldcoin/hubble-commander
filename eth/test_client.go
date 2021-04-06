package eth

import (
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
)

type TestClient struct {
	*Client
	*simulator.Simulator
}

func NewTestClient() (*TestClient, error) {
	sim, err := simulator.NewAutominingSimulator()
	if err != nil {
		return nil, err
	}

	contracts, err := deployer.DeployRollup(sim)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(sim, &NewClientParams{
		Rollup:          contracts.Rollup,
		AccountRegistry: contracts.AccountRegistry,
	})
	if err != nil {
		return nil, err
	}

	return &TestClient{
		Client:    client,
		Simulator: sim,
	}, nil
}
