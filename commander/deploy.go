package commander

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	// "github.com/Worldcoin/hubble-commander/eth"
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
		chooserAddress, _, err = deployer.DeployProofOfAuthority(
			blockchain,
			cfg.Ethereum.MineTimeout,
			[]common.Address{blockchain.GetAccount().From},
		)
		if err != nil {
			return nil, err
		}
	}

	totalGenesisAmount := models.NewUint256(0)
	for _, account := range cfg.Bootstrap.GenesisAccounts {
		totalGenesisAmount = totalGenesisAmount.Add(&account.State.Balance)
	}

	// Here we process the set of GenesisAccounts and build the tree we will deploy.

	accountTree := deployer.NewTree(st.AccountTreeDepth)
	for _, account := range cfg.Bootstrap.GenesisAccounts {
		accountTree.RegisterAccount(&account.PublicKey)
	}

	accountTreeRoot := accountTree.LeftRoot()
	accountSubtreesArray := (*[st.AccountTreeDepth-1]common.Hash)(accountTree.Subtrees)

	log.Printf(
		"- Using precomputed account tree.\n - root: %s\n - count: %d\n - subtrees: %s\n",
		accountTreeRoot,
		accountTree.LeafIndexLeft,
		accountSubtreesArray,
	)

	// lithp-PR: change signature, don't return accountRegistry
	accountRegistryAddress, accountRegistryDeploymentBlock, _, err := deployer.DeployAccountRegistry(
		blockchain,
		chooserAddress,
		cfg.Ethereum.MineTimeout,
		&accountTreeRoot,
		accountTree.LeafIndexLeft,
		*accountSubtreesArray,
	)
	if err != nil {
		return nil, err
	}

	// lithp-TODO: what does this do?
	//             This appears to build the eth object representing the registry
	//             I don't think it takes any on-chain actions.
	//             Probably doesn't require any changes.
	/*
	accountManager, err := eth.NewAccountManager(blockchain, &eth.AccountManagerParams{
		AccountRegistry:        accountRegistry,
		AccountRegistryAddress: *accountRegistryAddress,
	})
	if err != nil {
		return nil, err
	}
	*/

	// lithp-TODO: This saves the accounts to storage, it does not make on-chain changes
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
