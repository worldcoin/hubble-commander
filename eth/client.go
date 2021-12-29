package eth

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/metrics"
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
	SpokeRegistry   *spokeregistry.SpokeRegistry
	DepositManager  *depositmanager.DepositManager
	ClientConfig
}

type ClientConfig struct {
	TxTimeout                            *time.Duration  // default 60s
	StakeAmount                          *models.Uint256 // default 0.1 ether
	TransferBatchSubmissionGasLimit      *uint64         // default 400_000 gas
	C2TBatchSubmissionGasLimit           *uint64         // default 500_000 gas
	MassMigrationBatchSubmissionGasLimit *uint64         // default 550_000 gas
	DepositBatchSubmissionGasLimit       *uint64         // default 220_000 gas
	TransitionDisputeGasLimit            *uint64         // default 5_000_000 gas
	SignatureDisputeGasLimit             *uint64         // default 7_500_000 gas
	BatchAccountRegistrationGasLimit     *uint64         // default 8_000_000 gas
}

type Client struct {
	config                 ClientConfig
	ChainState             models.ChainState
	Blockchain             chain.Connection
	Metrics                *metrics.CommanderMetrics
	Rollup                 *Rollup
	TokenRegistry          *TokenRegistry
	SpokeRegistry          *SpokeRegistry
	DepositManager         *DepositManager
	blocksToFinalise       *int64
	maxDepositSubTreeDepth *uint8
	domain                 *bls.Domain

	*AccountManager
}

//goland:noinspection GoDeprecation
func NewClient(blockchain chain.Connection, commanderMetrics *metrics.CommanderMetrics, params *NewClientParams) (*Client, error) {
	fillWithDefaults(&params.ClientConfig)

	rollupAbi, err := abi.JSON(strings.NewReader(rollup.RollupABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tokenRegistryAbi, err := abi.JSON(strings.NewReader(tokenregistry.TokenRegistryABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	spokeRegistryAbi, err := abi.JSON(strings.NewReader(spokeregistry.SpokeRegistryABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	depositManagerAbi, err := abi.JSON(strings.NewReader(depositmanager.DepositManagerABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	backend := blockchain.GetBackend()
	accountManager, err := NewAccountManager(blockchain, &AccountManagerParams{
		AccountRegistry:                  params.AccountRegistry,
		AccountRegistryAddress:           params.ChainState.AccountRegistry,
		BatchAccountRegistrationGasLimit: params.BatchAccountRegistrationGasLimit,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	rollupContract := bind.NewBoundContract(params.ChainState.Rollup, rollupAbi, backend, backend, backend)
	tokenRegistryContract := bind.NewBoundContract(params.ChainState.TokenRegistry, tokenRegistryAbi, backend, backend, backend)
	spokeRegistryContract := bind.NewBoundContract(params.ChainState.SpokeRegistry, spokeRegistryAbi, backend, backend, backend)
	depositManagerContract := bind.NewBoundContract(params.ChainState.DepositManager, depositManagerAbi, backend, backend, backend)
	return &Client{
		config:         params.ClientConfig,
		ChainState:     params.ChainState,
		Blockchain:     blockchain,
		Metrics:        commanderMetrics,
		AccountManager: accountManager,
		Rollup: &Rollup{
			Rollup:   params.Rollup,
			Contract: MakeContract(&rollupAbi, rollupContract),
		},
		TokenRegistry: &TokenRegistry{
			TokenRegistry: params.TokenRegistry,
			Contract:      MakeContract(&tokenRegistryAbi, tokenRegistryContract),
		},
		SpokeRegistry: &SpokeRegistry{
			SpokeRegistry: params.SpokeRegistry,
			Contract:      MakeContract(&spokeRegistryAbi, spokeRegistryContract),
		},
		DepositManager: &DepositManager{
			DepositManager: params.DepositManager,
			Contract:       MakeContract(&depositManagerAbi, depositManagerContract),
		},
	}, nil
}

func fillWithDefaults(c *ClientConfig) {
	if c.TxTimeout == nil {
		c.TxTimeout = ref.Duration(60 * time.Second)
	}
	if c.StakeAmount == nil {
		c.StakeAmount = models.NewUint256(1e17)
	}
	if c.TransferBatchSubmissionGasLimit == nil {
		c.TransferBatchSubmissionGasLimit = ref.Uint64(config.DefaultTransferBatchSubmissionGasLimit)
	}
	if c.C2TBatchSubmissionGasLimit == nil {
		c.C2TBatchSubmissionGasLimit = ref.Uint64(config.DefaultC2TBatchSubmissionGasLimit)
	}
	if c.MassMigrationBatchSubmissionGasLimit == nil {
		c.MassMigrationBatchSubmissionGasLimit = ref.Uint64(config.DefaultMassMigrationBatchSubmissionGasLimit)
	}
	if c.DepositBatchSubmissionGasLimit == nil {
		c.DepositBatchSubmissionGasLimit = ref.Uint64(config.DefaultDepositBatchSubmissionGasLimit)
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
