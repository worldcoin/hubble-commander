package eth

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/pkg/errors"
)

type NewClientParams struct {
	ChainState      models.ChainState
	Rollup          *rollup.Rollup
	AccountRegistry *accountregistry.AccountRegistry
	ClientConfig
}

type ClientConfig struct {
	TxTimeout   *time.Duration  // default 60s
	StakeAmount *models.Uint256 // default 0.1 ether
}

type Client struct {
	config             ClientConfig
	ChainState         models.ChainState
	ChainConnection    deployer.ChainConnection
	Rollup             *rollup.Rollup
	RollupABI          *abi.ABI
	AccountRegistry    *accountregistry.AccountRegistry
	AccountRegistryABI *abi.ABI
	boundContract      *bind.BoundContract
	blocksToFinalise   *int64
}

func NewClient(chainConnection deployer.ChainConnection, params *NewClientParams) (*Client, error) {
	fillWithDefaults(&params.ClientConfig)

	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	accountRegistryAbi, err := abi.JSON(strings.NewReader(accountregistry.AccountRegistryABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	backend := chainConnection.GetBackend()
	boundContract := bind.NewBoundContract(params.ChainState.Rollup, rollupAbi, backend, backend, backend)
	return &Client{
		config:             params.ClientConfig,
		ChainState:         params.ChainState,
		ChainConnection:    chainConnection,
		Rollup:             params.Rollup,
		RollupABI:          &rollupAbi,
		AccountRegistry:    params.AccountRegistry,
		AccountRegistryABI: &accountRegistryAbi,
		boundContract:      boundContract,
	}, nil
}

func fillWithDefaults(c *ClientConfig) {
	if c.TxTimeout == nil {
		c.TxTimeout = ref.Duration(60 * time.Second)
	}
	if c.StakeAmount == nil {
		c.StakeAmount = models.NewUint256(1e17)
	}
}
