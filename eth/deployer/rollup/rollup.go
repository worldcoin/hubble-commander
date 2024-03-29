package rollup

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/libs/estimator"
	"github.com/Worldcoin/hubble-commander/contracts/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/proofofauthority"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/contracts/transfer"
	"github.com/Worldcoin/hubble-commander/contracts/vault"
	"github.com/Worldcoin/hubble-commander/contracts/withdrawmanager"
	"github.com/Worldcoin/hubble-commander/deployment"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultMaxDepositSubtreeDepth = 2
	DefaultGenesisStateRoot       = "cf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae"
	DefaultStakeAmount            = 1e15
	DefaultMinGasLeft             = 10_000
	DefaultMaxTxsPerCommit        = 32

	originalCostEstimatorAddress = "9f19c6456854fda109dfc2b9ad45d9bf147c8557"
)

type DeploymentConfig struct {
	Params
	Dependencies
	MineTimeout time.Duration
}

type Params struct {
	MaxDepositSubtreeDepth *models.Uint256
	GenesisStateRoot       *common.Hash
	StakeAmount            *models.Uint256
	BlocksToFinalise       *models.Uint256
	MinGasLeft             *models.Uint256
	MaxTxsPerCommit        *models.Uint256
	TotalGenesisAmount     *models.Uint256
}

type Dependencies struct {
	AccountRegistry *common.Address
	Chooser         *common.Address
}

type RollupContracts struct {
	Config                 DeploymentConfig
	Chooser                *proofofauthority.ProofOfAuthority
	AccountRegistry        *accountregistry.AccountRegistry
	AccountRegistryAddress common.Address
	TokenRegistry          *tokenregistry.TokenRegistry
	TokenRegistryAddress   common.Address
	SpokeRegistry          *spokeregistry.SpokeRegistry
	SpokeRegistryAddress   common.Address
	Vault                  *vault.Vault
	DepositManager         *depositmanager.DepositManager
	DepositManagerAddress  common.Address
	WithdrawManagerAddress common.Address
	Transfer               *transfer.Transfer
	MassMigration          *massmigration.MassMigration
	Create2Transfer        *create2transfer.Create2Transfer
	Rollup                 *rollup.Rollup
	RollupAddress          common.Address
	ExampleTokenAddress    common.Address
}

type txHelperContracts struct {
	TransferAddress        common.Address
	MassMigrationAddress   common.Address
	Create2TransferAddress common.Address
	TransferTx             *types.Transaction
	MassMigrationTx        *types.Transaction
	Create2TransferTx      *types.Transaction
	Transfer               *transfer.Transfer
	MassMigration          *massmigration.MassMigration
	Create2Transfer        *create2transfer.Create2Transfer
}

func DeployRollup(c chain.ManualNonceConnection) (*RollupContracts, error) {
	return DeployConfiguredRollup(c, &DeploymentConfig{})
}

// nolint:funlen,gocyclo
func DeployConfiguredRollup(c chain.ManualNonceConnection, cfg *DeploymentConfig) (*RollupContracts, error) {
	fillWithDefaults(&cfg.Params)

	waitForMultipleTxs := chain.CreateWaitForMultipleTxsHelper(c.GetBackend(), cfg.MineTimeout)

	// Stage 1
	err := deployMissing(&cfg.Dependencies, c, cfg.MineTimeout)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Stage 2
	log.Println("Deploying TokenRegistry")
	tokenRegistryAddress, tokenRegistryTx, tokenRegistry, err := tokenregistry.DeployTokenRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// It is possible that this method (and submethods) contains a bug.  We call
	// BumpNonce() whenever we successfully send a transaction.  We assume that an err
	// means a transaction was not sent and the nonce of the account does not need to
	// be incremented.  This might be a bad assumption.
	c.BumpNonce()

	log.Println("Deploying SpokeRegistry")
	spokeRegistryAddress, spokeRegistryTx, spokeRegistry, err := spokeregistry.DeploySpokeRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Deploying BNPairingPrecompileCostEstimator")
	costEstimatorAddress, costEstimatorDeployTx, costEstimator, err := estimator.DeployCostEstimator(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	_, err = waitForMultipleTxs(*tokenRegistryTx, *spokeRegistryTx, *costEstimatorDeployTx)
	if err != nil {
		return nil, err
	}

	// Stage 3
	log.Println("Initializing BNPairingPrecompileCostEstimator")
	costEstimatorInitTx, err := costEstimator.Run(c.GetAccount())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Deploying Vault")
	vaultAddress, vaultTx, vaultContract, err := vault.DeployVault(
		c.GetAccount(),
		c.GetBackend(),
		tokenRegistryAddress,
		spokeRegistryAddress,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	_, err = waitForMultipleTxs(*costEstimatorInitTx, *vaultTx)
	if err != nil {
		return nil, err
	}

	// Stage 4
	log.Println("Deploying DepositManager")
	depositManagerAddress, depositManagerTx, depositManager, err := depositmanager.DeployDepositManager(
		c.GetAccount(),
		c.GetBackend(),
		tokenRegistryAddress,
		vaultAddress,
		cfg.MaxDepositSubtreeDepth.ToBig(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	var txHelpers *txHelperContracts
	withReplacedCostEstimatorAddress(costEstimatorAddress, func() {
		txHelpers, err = deployTransactionHelperContracts(c)
	})
	if err != nil {
		return nil, err
	}

	_, err = waitForMultipleTxs(
		*depositManagerTx,
		*txHelpers.TransferTx,
		*txHelpers.MassMigrationTx,
		*txHelpers.Create2TransferTx,
	)
	if err != nil {
		return nil, err
	}

	// Stage 5
	log.Println("Deploying Rollup")
	stateRoot := [32]byte{}
	copy(stateRoot[:], cfg.GenesisStateRoot.Bytes())
	rollupAddress, tx, rollupContract, err := rollup.DeployRollup(
		c.GetAccount(),
		c.GetBackend(),
		*cfg.Chooser,
		depositManagerAddress,
		*cfg.AccountRegistry,
		txHelpers.TransferAddress,
		txHelpers.MassMigrationAddress,
		txHelpers.Create2TransferAddress,
		stateRoot,
		cfg.StakeAmount.ToBig(),
		cfg.BlocksToFinalise.ToBig(),
		cfg.MinGasLeft.ToBig(),
		cfg.MaxTxsPerCommit.ToBig(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	_, err = chain.WaitToBeMined(c.GetBackend(), cfg.MineTimeout, tx)
	if err != nil {
		return nil, err
	}

	// Stage 6
	log.Println("Setting Rollup address in DepositManager")
	depositManagerInitTx, err := depositManager.SetRollupAddress(c.GetAccount(), rollupAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Setting Rollup address in Vault")
	vaultInitTx, err := vaultContract.SetRollupAddress(c.GetAccount(), rollupAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Deploying WithdrawManager")
	var (
		withdrawManagerAddress common.Address
		withdrawManagerTx      *types.Transaction
	)
	withReplacedCostEstimatorAddress(costEstimatorAddress, func() {
		withdrawManagerAddress, withdrawManagerTx, _, err = withdrawmanager.DeployWithdrawManager(
			c.GetAccount(),
			c.GetBackend(),
			tokenRegistryAddress,
			vaultAddress,
			rollupAddress,
		)
	})
	c.BumpNonce()

	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Deploying TestCustomToken")
	exampleTokenAddress, exampleTokenTx, exampleToken, err := customtoken.DeployTestCustomToken(
		c.GetAccount(),
		c.GetBackend(),
		"ExampleToken",
		"EXP",
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	_, err = waitForMultipleTxs(
		*depositManagerInitTx,
		*vaultInitTx,
		*withdrawManagerTx,
		*exampleTokenTx,
	)
	if err != nil {
		return nil, err
	}

	// Stage 7
	stageSevenTxs := make([]types.Transaction, 0, 3)

	log.Println("Registering WithdrawManager as a spoke in SpokeRegistry")
	spokeRegistrationTx, err := spokeRegistry.RegisterSpoke(c.GetAccount(), withdrawManagerAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stageSevenTxs = append(stageSevenTxs, *spokeRegistrationTx)
	c.BumpNonce()

	log.Println("Registering TestCustomToken in TokenRegistry")
	registerTokenTx, err := tokenRegistry.RegisterToken(c.GetAccount(), exampleTokenAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	stageSevenTxs = append(stageSevenTxs, *registerTokenTx)
	c.BumpNonce()

	if cfg.TotalGenesisAmount != nil {
		log.Println("Transferring genesis funds to vault")
		var transferGenesisFundsTx *types.Transaction
		transferGenesisFundsTx, err = transferGenesisFunds(c, exampleToken, vaultAddress, cfg.TotalGenesisAmount)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		stageSevenTxs = append(stageSevenTxs, *transferGenesisFundsTx)
		c.BumpNonce()
	}

	_, err = waitForMultipleTxs(stageSevenTxs...)
	if err != nil {
		return nil, err
	}

	proofOfAuthority, err := proofofauthority.NewProofOfAuthority(*cfg.Chooser, c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(*cfg.AccountRegistry, c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &RollupContracts{
		Config:                 *cfg,
		Chooser:                proofOfAuthority,
		AccountRegistry:        accountRegistry,
		AccountRegistryAddress: *cfg.AccountRegistry,
		TokenRegistry:          tokenRegistry,
		TokenRegistryAddress:   tokenRegistryAddress,
		SpokeRegistry:          spokeRegistry,
		SpokeRegistryAddress:   spokeRegistryAddress,
		Vault:                  vaultContract,
		DepositManager:         depositManager,
		DepositManagerAddress:  depositManagerAddress,
		WithdrawManagerAddress: withdrawManagerAddress,
		Transfer:               txHelpers.Transfer,
		MassMigration:          txHelpers.MassMigration,
		Create2Transfer:        txHelpers.Create2Transfer,
		Rollup:                 rollupContract,
		RollupAddress:          rollupAddress,
		ExampleTokenAddress:    exampleTokenAddress,
	}, nil
}

func deployTransactionHelperContracts(c chain.ManualNonceConnection) (*txHelperContracts, error) {
	log.Println("Deploying Transfer")
	transferAddress, transferTx, transferContract, err := transfer.DeployTransfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Deploying MassMigration")
	massMigrationAddress, massMigrationTx, massMigration, err := massmigration.DeployMassMigration(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	log.Println("Deploying Create2Transfer")
	create2TransferAddress, create2TransferTx, create2Transfer, err := create2transfer.DeployCreate2Transfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}
	c.BumpNonce()

	return &txHelperContracts{
		TransferAddress:        transferAddress,
		MassMigrationAddress:   massMigrationAddress,
		Create2TransferAddress: create2TransferAddress,
		TransferTx:             transferTx,
		MassMigrationTx:        massMigrationTx,
		Create2TransferTx:      create2TransferTx,
		Transfer:               transferContract,
		MassMigration:          massMigration,
		Create2Transfer:        create2Transfer,
	}, nil
}

func fillWithDefaults(params *Params) {
	if params.MaxDepositSubtreeDepth == nil {
		params.MaxDepositSubtreeDepth = models.NewUint256(DefaultMaxDepositSubtreeDepth)
	}
	if params.GenesisStateRoot == nil {
		// Result of getDefaultGenesisRoot function from deploy.ts
		hash := common.HexToHash(DefaultGenesisStateRoot)
		params.GenesisStateRoot = &hash
	}
	if params.StakeAmount == nil {
		params.StakeAmount = models.NewUint256(DefaultStakeAmount)
	}
	if params.BlocksToFinalise == nil {
		params.BlocksToFinalise = models.NewUint256(uint64(config.DefaultBlocksToFinalise))
	}
	if params.MinGasLeft == nil {
		params.MinGasLeft = models.NewUint256(DefaultMinGasLeft)
	}
	if params.MaxTxsPerCommit == nil {
		params.MaxTxsPerCommit = models.NewUint256(DefaultMaxTxsPerCommit)
	}
}

func deployMissing(dependencies *Dependencies, c chain.ManualNonceConnection, mineTimeout time.Duration) error {
	if dependencies.Chooser == nil {
		proofOfAuthorityAddress, _, err := deployer.DeployProofOfAuthority(c, mineTimeout, []common.Address{c.GetAccount().From})
		if err != nil {
			return err
		}
		dependencies.Chooser = proofOfAuthorityAddress
		c.BumpNonce()
	}
	if dependencies.AccountRegistry == nil {
		tree := deployment.NewTree(storage.AccountTreeDepth)
		root := tree.LeftRoot()

		subtreesArray := (*[storage.AccountTreeDepth - 1]common.Hash)(tree.Subtrees)

		accountRegistryAddress, _, _, err := deployer.DeployAccountRegistry(c, dependencies.Chooser, mineTimeout, &root, 0, subtreesArray)
		if err != nil {
			return err
		}
		dependencies.AccountRegistry = accountRegistryAddress
		c.BumpNonce()
	}
	return nil
}

//goland:noinspection GoDeprecation
func withReplacedCostEstimatorAddress(newCostEstimator common.Address, fn func()) {
	targetString := strings.ToLower(newCostEstimator.String()[2:])
	originalWithdrawManagerBin := withdrawmanager.WithdrawManagerBin
	originalTransferBin := transfer.TransferBin
	originalCreate2TransferBin := create2transfer.Create2TransferBin
	originalMassMigrationBin := massmigration.MassMigrationBin

	withdrawmanager.WithdrawManagerBin = strings.Replace(originalWithdrawManagerBin, originalCostEstimatorAddress, targetString, -1)
	transfer.TransferBin = strings.Replace(originalTransferBin, originalCostEstimatorAddress, targetString, -1)
	create2transfer.Create2TransferBin = strings.Replace(originalCreate2TransferBin, originalCostEstimatorAddress, targetString, -1)
	massmigration.MassMigrationBin = strings.Replace(originalMassMigrationBin, originalCostEstimatorAddress, targetString, -1)

	defer func() {
		withdrawmanager.WithdrawManagerBin = originalWithdrawManagerBin
		transfer.TransferBin = originalTransferBin
		create2transfer.Create2TransferBin = originalCreate2TransferBin
		massmigration.MassMigrationBin = originalMassMigrationBin
	}()

	fn()
}

func transferGenesisFunds(
	c chain.Connection,
	tokenContract *customtoken.TestCustomToken,
	recipient common.Address,
	l2Amount *models.Uint256,
) (*types.Transaction, error) {
	return tokenContract.Transfer(c.GetAccount(), recipient, l2Amount.MulN(consts.L2Unit).ToBig())
}
