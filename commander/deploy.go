package commander

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func Deploy(cfg *config.DeployerConfig, blockchain chain.Connection) (chainSpec *string, err error) {
	tempStorage, err := st.NewTemporaryStorage()
	if err != nil {
		return nil, err
	}

	defer func() {
		closeErr := tempStorage.Close()
		if closeErr != nil {
			err = fmt.Errorf("temporary storage close caused by: %w, failed with: %v", err, closeErr)
		}
	}()

	log.Printf(
		"Bootstrapping genesis state with %d accounts on chainId = %d",
		len(cfg.Bootstrap.GenesisAccounts),
		cfg.Ethereum.ChainID,
	)
	chainState, err := deployContractsAndSetupGenesisState(tempStorage.Storage, blockchain, cfg)
	if err != nil {
		return nil, err
	}

	chainSpec, err = GenerateChainSpec(chainState)
	if err != nil {
		return nil, err
	}

	return chainSpec, nil
}

func deployContractsAndSetupGenesisState(
	storage *st.Storage,
	blockchain chain.Connection,
	cfg *config.DeployerConfig,
) (chainState *models.ChainState, err error) {
	var chooserAddress *common.Address
	if cfg.Bootstrap.Chooser != nil {
		chooserAddress = cfg.Bootstrap.Chooser
	} else {
		chooserAddress, _, err = deployer.DeployProofOfAuthority(blockchain, cfg.Ethereum.MineTimeout)
		if err != nil {
			return nil, err
		}
	}

	accountRegistryAddress, accountRegistryDeploymentBlock, accountRegistry, err := deployer.DeployAccountRegistry(
		blockchain,
		chooserAddress,
		cfg.Ethereum.MineTimeout,
	)
	if err != nil {
		return nil, err
	}

	accountManager, err := eth.NewAccountManager(blockchain, &eth.AccountManagerParams{
		AccountRegistry:        accountRegistry,
		AccountRegistryAddress: *accountRegistryAddress,
	})
	if err != nil {
		return nil, err
	}

	totalGenesisAmount, err := RegisterGenesisAccountsAndCalculateTotalAmount(
		accountManager,
		cfg.Bootstrap.GenesisAccounts,
		cfg.Ethereum.MineTimeout,
	)
	if err != nil {
		return nil, err
	}

	err = PopulateGenesisAccounts(storage, cfg.Bootstrap.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	stateRoot, err := storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	contracts, err := rollup.DeployConfiguredRollup(blockchain, &rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot:   stateRoot,
			BlocksToFinalise:   models.NewUint256(uint64(cfg.Bootstrap.BlocksToFinalise)),
			TotalGenesisAmount: totalGenesisAmount,
		},
		Dependencies: rollup.Dependencies{
			AccountRegistry: accountRegistryAddress,
			Chooser:         chooserAddress,
		},
		MineTimeout: cfg.Ethereum.MineTimeout,
	})
	if err != nil {
		return nil, err
	}

	chainState = &models.ChainState{
		ChainID:                        blockchain.GetChainID(),
		AccountRegistry:                *accountRegistryAddress,
		AccountRegistryDeploymentBlock: *accountRegistryDeploymentBlock,
		TokenRegistry:                  contracts.TokenRegistryAddress,
		SpokeRegistry:                  contracts.SpokeRegistryAddress,
		DepositManager:                 contracts.DepositManagerAddress,
		WithdrawManager:                contracts.WithdrawManagerAddress,
		Rollup:                         contracts.RollupAddress,
		GenesisAccounts:                cfg.Bootstrap.GenesisAccounts,
		SyncedBlock:                    getInitialSyncedBlock(*accountRegistryDeploymentBlock),
	}

	return chainState, nil
}
