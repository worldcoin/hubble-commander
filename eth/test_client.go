package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
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

	if cfg.Dependencies.AccountRegistry == nil {
		accountRegistryAddress, _, _, deployErr := accountregistry.DeployAccountRegistry(sim.GetAccount(), sim.GetBackend())
		if deployErr != nil {
			return nil, deployErr
		}

		cfg.Dependencies.AccountRegistry = &accountRegistryAddress
	}

	contracts, err := rollup.DeployConfiguredRollup(sim, cfg)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(sim, &NewClientParams{
		ChainState: models.ChainState{
			ChainID:         sim.GetChainID(),
			AccountRegistry: *cfg.AccountRegistry,
			DeploymentBlock: 0,
			TokenRegistry:   contracts.TokenRegistryAddress,
			DepositManager:  contracts.DepositManagerAddress,
			Rollup:          contracts.RollupAddress,
			SyncedBlock:     0,
			GenesisAccounts: nil,
		},
		Rollup:          contracts.Rollup,
		AccountRegistry: contracts.AccountRegistry,
		TokenRegistry:   contracts.TokenRegistry,
		DepositManager:  contracts.DepositManager,
		ClientConfig:    clientCfg,
	})
	if err != nil {
		return nil, err
	}

	return &TestClient{
		Client:    client,
		Simulator: sim,
	}, nil
}
