package eth

import (
	"context"
	"sync"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type TestClient struct {
	*Client
	*simulator.Simulator
	ExampleTokenAddress common.Address

	cancelTxsSending context.CancelFunc
	wg               sync.WaitGroup
	TxsChannels      *TxsTrackingChannels
}

type TestClientConfig struct {
	TxsChannels *TxsTrackingChannels
	ClientConfig
}

var (
	DomainOnlyTestClient = &Client{domain: &bls.TestDomain}
)

// NewTestClient Sets up a TestClient backed by automining simulator.
// Remember to call Close() at the end of the test
func NewTestClient() (*TestClient, error) {
	return NewConfiguredTestClient(&rollup.DeploymentConfig{}, &TestClientConfig{})
}

func NewConfiguredTestClient(cfg *rollup.DeploymentConfig, clientCfg *TestClientConfig) (*TestClient, error) {
	sim, err := simulator.NewAutominingSimulator()
	if err != nil {
		return nil, err
	}

	contracts, err := rollup.DeployConfiguredRollup(sim, cfg)
	if err != nil {
		return nil, err
	}

	startTxsSending := clientCfg.TxsChannels == nil
	if startTxsSending {
		clientCfg.TxsChannels = &TxsTrackingChannels{
			Requests: make(chan *TxSendingRequest, 32),
			SentTxs:  make(chan *types.Transaction, 32),
		}
	}

	client, err := NewClient(sim, metrics.NewCommanderMetrics(), &NewClientParams{
		ChainState: models.ChainState{
			ChainID:                        sim.GetChainID(),
			AccountRegistry:                contracts.AccountRegistryAddress,
			AccountRegistryDeploymentBlock: 0,
			TokenRegistry:                  contracts.TokenRegistryAddress,
			SpokeRegistry:                  contracts.SpokeRegistryAddress,
			DepositManager:                 contracts.DepositManagerAddress,
			Rollup:                         contracts.RollupAddress,
			SyncedBlock:                    0,
			GenesisAccounts:                nil,
		},
		Rollup:          contracts.Rollup,
		AccountRegistry: contracts.AccountRegistry,
		TokenRegistry:   contracts.TokenRegistry,
		SpokeRegistry:   contracts.SpokeRegistry,
		DepositManager:  contracts.DepositManager,
		ClientConfig:    clientCfg.ClientConfig,
		TxsChannels:     clientCfg.TxsChannels,
	})
	if err != nil {
		return nil, err
	}

	testClient := &TestClient{
		Client:              client,
		Simulator:           sim,
		ExampleTokenAddress: contracts.ExampleTokenAddress,
		TxsChannels:         clientCfg.TxsChannels,
		cancelTxsSending:    func() {},
	}

	if startTxsSending {
		testClient.startTxsSending()
	}
	return testClient, nil
}

func (c *TestClient) Close() {
	c.stopTxsSending()
	c.Simulator.Close()
}

func (c *TestClient) startTxsSending() {
	ctx, cancel := context.WithCancel(context.Background())
	c.cancelTxsSending = cancel

	c.wg.Add(1)
	go func() {
		err := c.sendTxs(ctx, c.TxsChannels.Requests)
		if err != nil {
			panic(err)
		}
		c.wg.Done()
	}()
}

func (c *TestClient) stopTxsSending() {
	c.cancelTxsSending()
	c.wg.Wait()
}

func (c *TestClient) sendTxs(ctx context.Context, requests <-chan *TxSendingRequest) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case req := <-requests:
			err := req.Send()
			if err != nil {
				return err
			}
		}
	}
}
