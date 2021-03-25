package deployer

import (
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
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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
	Config          DeploymentConfig
	Chooser         *proofofburn.ProofOfBurn
	AccountRegistry *accountregistry.AccountRegistry
	TokenRegistry   *tokenregistry.TokenRegistry
	SpokeRegistry   *spokeregistry.SpokeRegistry
	Vault           *vault.Vault
	DepositManager  *depositmanager.DepositManager
	Transfer        *transfer.Transfer
	MassMigration   *massmigration.MassMigration
	Create2Transfer *create2transfer.Create2Transfer
	Rollup          *rollup.Rollup
	RollupAddress   common.Address
}

func DeployRollup(d ChainConnection) (*RollupContracts, error) {
	accountRegistryAddress, _, err := DeployAccountRegistry(d)
	if err != nil {
		return nil, err
	}
	return DeployConfiguredRollup(d, DeploymentConfig{
		AccountRegistryAddress: accountRegistryAddress,
	})
}

func DeployConfiguredRollup(d ChainConnection, config DeploymentConfig) (*RollupContracts, error) {
	fillWithDefaults(&config)
	proofOfBurnAddress, _, proofOfBurn, err := proofofburn.DeployProofOfBurn(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	tokenRegistryAddress, _, tokenRegistry, err := tokenregistry.DeployTokenRegistry(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	spokeRegistryAddress, _, spokeRegistry, err := spokeregistry.DeploySpokeRegistry(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	vaultAddress, _, vaultContract, err := vault.DeployVault(
		d.TransactionOpts(),
		d.GetBackend(),
		tokenRegistryAddress,
		spokeRegistryAddress,
	)
	if err != nil {
		return nil, err
	}

	depositManagerAddress, _, depositManager, err := depositmanager.DeployDepositManager(
		d.TransactionOpts(),
		d.GetBackend(),
		tokenRegistryAddress,
		vaultAddress,
		&config.MaxDepositSubtreeDepth.Int,
	)
	if err != nil {
		return nil, err
	}

	accountRegistry, err := accountregistry.NewAccountRegistry(*config.AccountRegistryAddress, d.GetBackend())
	if err != nil {
		return nil, err
	}

	transferAddress, _, transferContract, err := transfer.DeployTransfer(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	massMigrationAddress, _, massMigration, err := massmigration.DeployMassMigration(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	d.Commit()

	create2TransferAddress, _, create2Transfer, err := create2transfer.DeployCreate2Transfer(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	stateRoot := [32]byte{}
	copy(stateRoot[:], config.GenesisStateRoot.Bytes())
	rollupAddress, _, rollupContract, err := rollup.DeployRollup(
		d.TransactionOpts(),
		d.GetBackend(),
		proofOfBurnAddress,
		depositManagerAddress,
		*config.AccountRegistryAddress,
		transferAddress,
		massMigrationAddress,
		create2TransferAddress,
		stateRoot,
		&config.StakeAmount.Int,
		&config.BlocksToFinalise.Int,
		&config.MinGasLeft.Int,
		&config.MaxTxsPerCommit.Int,
	)
	if err != nil {
		return nil, err
	}

	d.Commit()

	return &RollupContracts{
		Config:          config,
		Chooser:         proofOfBurn,
		AccountRegistry: accountRegistry,
		TokenRegistry:   tokenRegistry,
		SpokeRegistry:   spokeRegistry,
		Vault:           vaultContract,
		DepositManager:  depositManager,
		Transfer:        transferContract,
		MassMigration:   massMigration,
		Create2Transfer: create2Transfer,
		Rollup:          rollupContract,
		RollupAddress:   rollupAddress,
	}, nil
}

func fillWithDefaults(config *DeploymentConfig) {
	if config.MaxDepositSubtreeDepth == nil {
		config.MaxDepositSubtreeDepth = models.NewUint256(2)
	}
	if config.GenesisStateRoot == nil {
		// Result of getDefaultGenesisRoot function from deploy.ts
		hash := common.HexToHash("cf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae")
		config.GenesisStateRoot = &hash
	}
	if config.StakeAmount == nil {
		config.StakeAmount = models.NewUint256(1e17)
	}
	if config.BlocksToFinalise == nil {
		config.BlocksToFinalise = models.NewUint256(7 * 24 * 60 * 4)
	}
	if config.MinGasLeft == nil {
		config.MinGasLeft = models.NewUint256(10_000)
	}
	if config.MaxTxsPerCommit == nil {
		config.MaxTxsPerCommit = models.NewUint256(32)
	}
}
