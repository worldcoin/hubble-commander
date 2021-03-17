package eth

import (
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NewClientParams struct {
	ethNodeAddress         string
	rollupAddress          common.Address
	accountRegistryAddress common.Address
	ClientConfig
}

type ClientConfig struct {
	txTimeout   *time.Duration  // default 60s
	stakeAmount *models.Uint256 // default 0.1 ether
}

type Client struct {
	account         bind.TransactOpts
	config          ClientConfig
	Rollup          *rollup.Rollup
	AccountRegistry *accountregistry.AccountRegistry
}

func NewClient(account *bind.TransactOpts, params NewClientParams) (*Client, error) {
	fillWithDefaults(&params.ClientConfig)

	backend, err := ethclient.Dial(params.ethNodeAddress)
	if err != nil {
		return nil, err
	}

	rollupContract, err := rollup.NewRollup(params.rollupAddress, backend)
	if err != nil {
		return nil, err
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(params.accountRegistryAddress, backend)
	if err != nil {
		return nil, err
	}

	return &Client{
		account:         *account,
		config:          params.ClientConfig,
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
	}, nil
}

func NewTestClient(
	account *bind.TransactOpts,
	rollupContract *rollup.Rollup,
	accountRegistry *accountregistry.AccountRegistry,
) *Client {
	return &Client{
		account:         *account,
		config:          getDefaultConfig(),
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
	}
}

func fillWithDefaults(c *ClientConfig) {
	if c.txTimeout == nil {
		c.txTimeout = ref.Duration(60 * time.Second)
	}
	if c.stakeAmount == nil {
		c.stakeAmount = models.NewUint256(1e17)
	}
}

func getDefaultConfig() ClientConfig {
	var config ClientConfig
	fillWithDefaults(&config)
	return config
}
