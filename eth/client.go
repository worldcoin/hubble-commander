package eth

import (
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type NewClientParams struct {
	Rollup          *rollup.Rollup
	AccountRegistry *accountregistry.AccountRegistry
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

	return &Client{
		account:         *account,
		config:          params.ClientConfig,
		Rollup:          params.Rollup,
		AccountRegistry: params.AccountRegistry,
	}, nil
}

func fillWithDefaults(c *ClientConfig) {
	if c.txTimeout == nil {
		c.txTimeout = ref.Duration(60 * time.Second)
	}
	if c.stakeAmount == nil {
		c.stakeAmount = models.NewUint256(1e17)
	}
}
