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
	"github.com/Worldcoin/hubble-commander/eth/chain"
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
	TxTimeout                        *time.Duration  // default 60s
	StakeAmount                      *models.Uint256 // default 0.1 ether
	TransitionDisputeGasLimit        *uint64         // default 5_000_000 gas
	SignatureDisputeGasLimit         *uint64         // default 7_500_000 gas
	BatchAccountRegistrationGasLimit *uint64         // default 8_000_000 gas
}

type Client struct {
	config                 ClientConfig
	ChainState             models.ChainState
	Blockchain             chain.Connection
	Rollup                 *Rollup
	TokenRegistry          *TokenRegistry
	DepositManager         *depositmanager.DepositManager
	DepositManagerABI      *abi.ABI
	depositManagerContract *bind.BoundContract
	blocksToFinalise       *int64
	maxDepositSubTreeDepth *uint8
	domain                 *bls.Domain

	*AccountManager
}

//goland:noinspection GoDeprecation
func NewClient(blockchain chain.Connection, params *NewClientParams) (*Client, error) {
	fillWithDefaults(&params.ClientConfig)

	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenRegistryAbi, err := abi.JSON(strings.NewReader(tokenregistry.TokenRegistryABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	depositManagerAbi, err := abi.JSON(strings.NewReader(depositmanager.DepositManagerABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	backend := blockchain.GetBackend()
	rollupContract := bind.NewBoundContract(params.ChainState.Rollup, rollupAbi, backend, backend, backend)
	tokenRegistryContract := bind.NewBoundContract(params.ChainState.TokenRegistry, tokenRegistryAbi, backend, backend, backend)
	depositManagerContract := bind.NewBoundContract(params.ChainState.DepositManager, depositManagerAbi, backend, backend, backend)
	accountManager, err := NewAccountManager(blockchain, &AccountManagerParams{
		AccountRegistry:                  params.AccountRegistry,
		AccountRegistryAddress:           params.ChainState.AccountRegistry,
		BatchAccountRegistrationGasLimit: params.BatchAccountRegistrationGasLimit,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &Client{
		config:     params.ClientConfig,
		ChainState: params.ChainState,
		Blockchain: blockchain,
		Rollup: &Rollup{
			Rollup: params.Rollup,
			Contract: Contract{
				ABI:           &rollupAbi,
				BoundContract: rollupContract,
			},
		},
		TokenRegistry: &TokenRegistry{
			TokenRegistry: params.TokenRegistry,
			Contract: Contract{
				ABI:           &tokenRegistryAbi,
				BoundContract: tokenRegistryContract,
			},
		},
		DepositManager:         params.DepositManager,
		DepositManagerABI:      &depositManagerAbi,
		depositManagerContract: depositManagerContract,
		AccountManager:         accountManager,
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
	if c.BatchAccountRegistrationGasLimit == nil {
		c.BatchAccountRegistrationGasLimit = ref.Uint64(config.DefaultBatchAccountRegistrationGasLimit)
	}
}
