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
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
)

type DeploymentConfig struct {
	MaxDepositSubtreeDepth *models.Uint256
	GenesisStateRoot       *models.Bytes32
	StakeAmount            *models.Uint256
	BlocksToFinalise       *models.Uint256
	MinGasLeft             *models.Uint256
	MaxTxsPerCommit        *models.Uint256
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
}

func DeployRollup(sim *simulator.Simulator) (*RollupContracts, error) {
	return DeployConfiguredRollup(sim, DeploymentConfig{})
}

func DeployConfiguredRollup(sim *simulator.Simulator, config DeploymentConfig) (*RollupContracts, error) {
	fillWithDefaults(&config)
	deployer := sim.Account

	proofOfBurnAddress, _, proofOfBurn, err := proofofburn.DeployProofOfBurn(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	tokenRegistryAddress, _, tokenRegistry, err := tokenregistry.DeployTokenRegistry(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	spokeRegistryAddress, _, spokeRegistry, err := spokeregistry.DeploySpokeRegistry(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	vaultAddress, _, vaultContract, err := vault.DeployVault(
		deployer,
		sim.Backend,
		tokenRegistryAddress,
		spokeRegistryAddress,
	)
	if err != nil {
		return nil, err
	}

	depositManagerAddress, _, depositManager, err := depositmanager.DeployDepositManager(
		deployer,
		sim.Backend,
		tokenRegistryAddress,
		vaultAddress,
		&config.MaxDepositSubtreeDepth.Int,
	)
	if err != nil {
		return nil, err
	}

	accountRegistryAddress, _, accountRegistry, err := accountregistry.DeployAccountRegistry(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	transferAddress, _, transferContract, err := transfer.DeployTransfer(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	massMigrationAddress, _, massMigration, err := massmigration.DeployMassMigration(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	sim.Backend.Commit()

	create2TransferAddress, _, create2Transfer, err := create2transfer.DeployCreate2Transfer(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	_, _, rollupContract, err := rollup.DeployRollup(
		deployer,
		sim.Backend,
		proofOfBurnAddress,
		depositManagerAddress,
		accountRegistryAddress,
		transferAddress,
		massMigrationAddress,
		create2TransferAddress,
		config.GenesisStateRoot.Bytes,
		&config.StakeAmount.Int,
		&config.BlocksToFinalise.Int,
		&config.MinGasLeft.Int,
		&config.MaxTxsPerCommit.Int,
	)
	if err != nil {
		return nil, err
	}

	sim.Backend.Commit()

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
	}, nil
}

func fillWithDefaults(config *DeploymentConfig) {
	if config.MaxDepositSubtreeDepth == nil {
		config.MaxDepositSubtreeDepth = utils.Uint256(models.MakeUint256(2))
	}
	if config.GenesisStateRoot == nil {
		// Result of getDefaultGenesisRoot function from deploy.ts
		root, _ := models.MakeBytes32("cf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae")
		config.GenesisStateRoot = &root
	}
	if config.StakeAmount == nil {
		config.StakeAmount = utils.Uint256(models.MakeUint256(1e17))
	}
	if config.BlocksToFinalise == nil {
		config.BlocksToFinalise = utils.Uint256(models.MakeUint256(7 * 24 * 60 * 4))
	}
	if config.MinGasLeft == nil {
		config.MinGasLeft = utils.Uint256(models.MakeUint256(10_000))
	}
	if config.MaxTxsPerCommit == nil {
		config.MaxTxsPerCommit = utils.Uint256(models.MakeUint256(32))
	}
}
