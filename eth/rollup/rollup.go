package rollup

import (
	"context"
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/proofofburn"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/Worldcoin/hubble-commander/contracts/transfer"
	"github.com/Worldcoin/hubble-commander/contracts/vault"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const (
	DefaultMaxDepositSubtreeDepth = 2
	DefaultGenesisStateRoot       = "cf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae"
	DefaultStakeAmount            = 1e17
	DefaultBlocksToFinalise       = 7 * 24 * 60 * 4
	DefaultMinGasLeft             = 10_000
	DefaultMaxTxsPerCommit        = 32
)

type DeploymentConfig struct {
	MaxDepositSubtreeDepth *models.Uint256
	GenesisStateRoot       *common.Hash
	StakeAmount            *models.Uint256
	BlocksToFinalise       *models.Uint256
	MinGasLeft             *models.Uint256
	MaxTxsPerCommit        *models.Uint256
	AccountRegistryAddress *common.Address
}

type RollupContracts struct {
	Config                DeploymentConfig
	Chooser               *proofofburn.ProofOfBurn
	AccountRegistry       *accountregistry.AccountRegistry
	TokenRegistry         *tokenregistry.TokenRegistry
	SpokeRegistry         *spokeregistry.SpokeRegistry
	Vault                 *vault.Vault
	DepositManager        *depositmanager.DepositManager
	Transfer              *transfer.Transfer
	MassMigration         *massmigration.MassMigration
	Create2Transfer       *create2transfer.Create2Transfer
	Rollup                *rollup.Rollup
	RollupAddress         common.Address
	RollupDeploymentBlock uint64
}

func DeployRollup(c deployer.ChainConnection) (*RollupContracts, error) {
	accountRegistryAddress, _, err := deployer.DeployAccountRegistry(c)
	if err != nil {
		return nil, err
	}
	return DeployConfiguredRollup(c, DeploymentConfig{
		AccountRegistryAddress: accountRegistryAddress,
	})
}

// nolint:funlen,gocyclo
func DeployConfiguredRollup(c deployer.ChainConnection, config DeploymentConfig) (*RollupContracts, error) {
	fillWithDefaults(&config)

	log.Println("Deploying ProofOfBurn")
	proofOfBurnAddress, tx, proofOfBurn, err := proofofburn.DeployProofOfBurn(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying TokenRegistry")
	tokenRegistryAddress, tx, tokenRegistry, err := tokenregistry.DeployTokenRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying SpokeRegistry")
	spokeRegistryAddress, tx, spokeRegistry, err := spokeregistry.DeploySpokeRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
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

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying DepositManager")
	depositManagerAddress, tx, depositManager, err := depositmanager.DeployDepositManager(
		c.GetAccount(),
		c.GetBackend(),
		tokenRegistryAddress,
		vaultAddress,
		config.MaxDepositSubtreeDepth.ToBig(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(*config.AccountRegistryAddress, c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Println("Deploying Transfer")
	transferAddress, tx, transferContract, err := transfer.DeployTransfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying MassMigration")
	massMigrationAddress, tx, massMigration, err := massmigration.DeployMassMigration(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying Create2Transfer")
	create2TransferAddress, tx, create2Transfer, err := create2transfer.DeployCreate2Transfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	log.Println("Deploying Rollup")
	stateRoot := [32]byte{}
	copy(stateRoot[:], config.GenesisStateRoot.Bytes())
	rollupAddress, tx, rollupContract, err := rollup.DeployRollup(
		c.GetAccount(),
		c.GetBackend(),
		proofOfBurnAddress,
		depositManagerAddress,
		*config.AccountRegistryAddress,
		transferAddress,
		massMigrationAddress,
		create2TransferAddress,
		stateRoot,
		config.StakeAmount.ToBig(),
		config.BlocksToFinalise.ToBig(),
		config.MinGasLeft.ToBig(),
		config.MaxTxsPerCommit.ToBig(),
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	c.Commit()
	_, err = deployer.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, err
	}

	rollupTxReceipt, err := c.GetBackend().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, err
	}

	return &RollupContracts{
		Config:                config,
		Chooser:               proofOfBurn,
		AccountRegistry:       accountRegistry,
		TokenRegistry:         tokenRegistry,
		SpokeRegistry:         spokeRegistry,
		Vault:                 vaultContract,
		DepositManager:        depositManager,
		Transfer:              transferContract,
		MassMigration:         massMigration,
		Create2Transfer:       create2Transfer,
		Rollup:                rollupContract,
		RollupAddress:         rollupAddress,
		RollupDeploymentBlock: rollupTxReceipt.BlockNumber.Uint64(),
	}, nil
}

func fillWithDefaults(config *DeploymentConfig) {
	if config.MaxDepositSubtreeDepth == nil {
		config.MaxDepositSubtreeDepth = models.NewUint256(DefaultMaxDepositSubtreeDepth)
	}
	if config.GenesisStateRoot == nil {
		// Result of getDefaultGenesisRoot function from deploy.ts
		hash := common.HexToHash(DefaultGenesisStateRoot)
		config.GenesisStateRoot = &hash
	}
	if config.StakeAmount == nil {
		config.StakeAmount = models.NewUint256(DefaultStakeAmount)
	}
	if config.BlocksToFinalise == nil {
		config.BlocksToFinalise = models.NewUint256(DefaultBlocksToFinalise)
	}
	if config.MinGasLeft == nil {
		config.MinGasLeft = models.NewUint256(DefaultMinGasLeft)
	}
	if config.MaxTxsPerCommit == nil {
		config.MaxTxsPerCommit = models.NewUint256(DefaultMaxTxsPerCommit)
	}
}
