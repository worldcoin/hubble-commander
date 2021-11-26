package eth

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/common"
)

type TestClient struct {
	*Client
	*simulator.Simulator
	ExampleTokenAddress common.Address
}

var (
	DomainOnlyTestClient = &Client{domain: &bls.TestDomain}
)

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

	client, err := NewClient(sim, metrics.NewCommanderMetrics(), &NewClientParams{
		ChainState: models.ChainState{
			ChainID:                        sim.GetChainID(),
			AccountRegistry:                contracts.AccountRegistryAddress,
			AccountRegistryDeploymentBlock: 0,
			TokenRegistry:                  contracts.TokenRegistryAddress,
			DepositManager:                 contracts.DepositManagerAddress,
			Rollup:                         contracts.RollupAddress,
			SyncedBlock:                    0,
			GenesisAccounts:                nil,
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
		Client:              client,
		Simulator:           sim,
		ExampleTokenAddress: contracts.ExampleTokenAddress,
	}, nil
}
