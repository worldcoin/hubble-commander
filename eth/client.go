package eth

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
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
	TokenRegistry   *tokenregistry.TokenRegistry
	DepositManager  *depositmanager.DepositManager
	ClientConfig
}

type ClientConfig struct {
	TxTimeout                 *time.Duration  // default 60s
	StakeAmount               *models.Uint256 // default 0.1 ether
	TransitionDisputeGasLimit *uint64         // default 5_000_000 gas
	SignatureDisputeGasLimit  *uint64         // default 7_500_000 gas
}

type Client struct {
	config                  ClientConfig
	ChainState              models.ChainState
	ChainConnection         deployer.ChainConnection
	Rollup                  *rollup.Rollup
	RollupABI               *abi.ABI
	AccountRegistry         *accountregistry.AccountRegistry
	AccountRegistryABI      *abi.ABI
	TokenRegistry           *tokenregistry.TokenRegistry
	DepositManager          *depositmanager.DepositManager
	rollupContract          *bind.BoundContract
	accountRegistryContract *bind.BoundContract
	blocksToFinalise        *int64
	domain                  *bls.Domain
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
	rollupContract := bind.NewBoundContract(params.ChainState.Rollup, rollupAbi, backend, backend, backend)
	accountRegistryContract := bind.NewBoundContract(params.ChainState.AccountRegistry, accountRegistryAbi, backend, backend, backend)
	return &Client{
		config:                  params.ClientConfig,
		ChainState:              params.ChainState,
		ChainConnection:         chainConnection,
		Rollup:                  params.Rollup,
		RollupABI:               &rollupAbi,
		AccountRegistry:         params.AccountRegistry,
		AccountRegistryABI:      &accountRegistryAbi,
		TokenRegistry:           params.TokenRegistry,
		DepositManager:          params.DepositManager,
		rollupContract:          rollupContract,
		accountRegistryContract: accountRegistryContract,
	}, nil
}

func fillWithDefaults(c *ClientConfig) {
	if c.TxTimeout == nil {
		c.TxTimeout = ref.Duration(60 * time.Second)
	}
	if c.StakeAmount == nil {
		c.StakeAmount = models.NewUint256(1e17)
	}
	if c.TransitionDisputeGasLimit == nil {
		c.TransitionDisputeGasLimit = ref.Uint64(config.DefaultTransitionDisputeGasLimit)
	}
	if c.SignatureDisputeGasLimit == nil {
		c.SignatureDisputeGasLimit = ref.Uint64(config.DefaultSignatureDisputeGasLimit)
	}
}
