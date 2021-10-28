package rollup

import (
	"strings"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/libs/estimator"
	"github.com/Worldcoin/hubble-commander/contracts/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/proofofburn"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/test/exampletoken"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/contracts/transfer"
	"github.com/Worldcoin/hubble-commander/contracts/vault"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultMaxDepositSubtreeDepth = 2
	DefaultGenesisStateRoot       = "cf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae"
	DefaultStakeAmount            = 1e17
	DefaultMinGasLeft             = 10_000
	DefaultMaxTxsPerCommit        = 32

	costEstimatorAddress = "079d8077c465bd0bf0fc502ad2b846757e415661"
)

type DeploymentConfig struct {
	Params
	Dependencies
}

type Params struct {
	MaxDepositSubtreeDepth *models.Uint256
	GenesisStateRoot       *common.Hash
	StakeAmount            *models.Uint256
	BlocksToFinalise       *models.Uint256
	MinGasLeft             *models.Uint256
	MaxTxsPerCommit        *models.Uint256
}

type Dependencies struct {
	AccountRegistry *common.Address
	Chooser         *common.Address
}

type RollupContracts struct {
	Config                 DeploymentConfig
	Chooser                *proofofburn.ProofOfBurn
	AccountRegistry        *accountregistry.AccountRegistry
	AccountRegistryAddress common.Address
	TokenRegistry          *tokenregistry.TokenRegistry
	TokenRegistryAddress   common.Address
	SpokeRegistry          *spokeregistry.SpokeRegistry
	Vault                  *vault.Vault
	DepositManager         *depositmanager.DepositManager
	DepositManagerAddress  common.Address
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
	Transfer               *transfer.Transfer
	MassMigration          *massmigration.MassMigration
	Create2Transfer        *create2transfer.Create2Transfer
}

func DeployRollup(c chain.Connection) (*RollupContracts, error) {
	return DeployConfiguredRollup(c, DeploymentConfig{})
}

// nolint:funlen,gocyclo
func DeployConfiguredRollup(c chain.Connection, cfg DeploymentConfig) (*RollupContracts, error) {
	fillWithDefaults(&cfg.Params)

	err := deployMissing(&cfg.Dependencies, c)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Deploying TokenRegistry")
	tokenRegistryAddress, tx, tokenRegistry, err := tokenregistry.DeployTokenRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying SpokeRegistry")
	spokeRegistryAddress, tx, spokeRegistry, err := spokeregistry.DeploySpokeRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying Vault")
	vaultAddress, tx, vaultContract, err := vault.DeployVault(
		c.GetAccount(),
		c.GetBackend(),
		tokenRegistryAddress,
		spokeRegistryAddress,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying DepositManager")
	depositManagerAddress, tx, depositManager, err := depositmanager.DeployDepositManager(
		c.GetAccount(),
		c.GetBackend(),
		tokenRegistryAddress,
		vaultAddress,
		cfg.MaxDepositSubtreeDepth.ToBig(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	exampleTokenAddress, tx, _, err := exampletoken.DeployTestExampleToken(
		c.GetAccount(),
		c.GetBackend(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	proofOfBurn, err := proofofburn.NewProofOfBurn(*cfg.Chooser, c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(*cfg.AccountRegistry, c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	estimatorAddress, err := deployCostEstimator(c)
	if err != nil {
		return nil, err
	}

	var txHelpers *txHelperContracts
	withReplacedCostEstimatorAddress(*estimatorAddress, func() {
		txHelpers, err = deployTransactionHelperContracts(c)
	})
	if err != nil {
		return nil, err
	}

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

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	tx, err = depositManager.SetRollupAddress(c.GetAccount(), rollupAddress)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return &RollupContracts{
		Config:                 cfg,
		Chooser:                proofOfBurn,
		AccountRegistry:        accountRegistry,
		AccountRegistryAddress: *cfg.AccountRegistry,
		TokenRegistry:          tokenRegistry,
		TokenRegistryAddress:   tokenRegistryAddress,
		SpokeRegistry:          spokeRegistry,
		Vault:                  vaultContract,
		DepositManager:         depositManager,
		DepositManagerAddress:  depositManagerAddress,
		Transfer:               txHelpers.Transfer,
		MassMigration:          txHelpers.MassMigration,
		Create2Transfer:        txHelpers.Create2Transfer,
		Rollup:                 rollupContract,
		RollupAddress:          rollupAddress,
		ExampleTokenAddress:    exampleTokenAddress,
	}, nil
}

func deployCostEstimator(c chain.Connection) (*common.Address, error) {
	log.Println("Deploying BNPairingPrecompileCostEstimator")
	estimatorAddress, tx, costEstimator, err := estimator.DeployCostEstimator(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	_, err = costEstimator.Run(c.GetAccount())
	if err != nil {
		return nil, err
	}

	return &estimatorAddress, nil
}

func deployTransactionHelperContracts(c chain.Connection) (*txHelperContracts, error) {
	log.Println("Deploying Transfer")
	transferAddress, tx, transferContract, err := transfer.DeployTransfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying MassMigration")
	massMigrationAddress, tx, massMigration, err := massmigration.DeployMassMigration(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying Create2Transfer")
	create2TransferAddress, tx, create2Transfer, err := create2transfer.DeployCreate2Transfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	_, err = chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	return &txHelperContracts{
		TransferAddress:        transferAddress,
		MassMigrationAddress:   massMigrationAddress,
		Create2TransferAddress: create2TransferAddress,
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

func deployMissing(dependencies *Dependencies, c chain.Connection) error {
	if dependencies.Chooser == nil {
		proofOfBurnAddress, _, err := deployer.DeployProofOfBurn(c)
		if err != nil {
			return err
		}
		dependencies.Chooser = proofOfBurnAddress
	}
	if dependencies.AccountRegistry == nil {
		accountRegistryAddress, _, _, err := deployer.DeployAccountRegistry(c, dependencies.Chooser)
		if err != nil {
			return err
		}
		dependencies.AccountRegistry = accountRegistryAddress
	}
	return nil
}

//goland:noinspection GoDeprecation
func withReplacedCostEstimatorAddress(newCostEstimator common.Address, fn func()) {
	targetString := strings.ToLower(newCostEstimator.String()[2:])
	originalTransferBin := transfer.TransferBin
	originalCreate2TransferBin := create2transfer.Create2TransferBin
	originalMassMigrationBin := massmigration.MassMigrationBin

	transfer.TransferBin = strings.Replace(originalTransferBin, costEstimatorAddress, targetString, -1)
	create2transfer.Create2TransferBin = strings.Replace(originalCreate2TransferBin, costEstimatorAddress, targetString, -1)
	massmigration.MassMigrationBin = strings.Replace(originalMassMigrationBin, costEstimatorAddress, targetString, -1)

	defer func() {
		transfer.TransferBin = originalTransferBin
		create2transfer.Create2TransferBin = originalCreate2TransferBin
		massmigration.MassMigrationBin = originalMassMigrationBin
	}()

	fn()
}
