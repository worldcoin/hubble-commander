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
	log "github.com/sirupsen/logrus"
)

var ErrNoPublicKeysInGenesisAccounts = fmt.Errorf("genesis accounts for deployment require public keys")

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
	chainState, err := deployContractsAndSetupGenesisState(tempStorage.Storage, blockchain, cfg.Bootstrap)
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
	cfg *config.DeployerBootstrapConfig,
) (*models.ChainState, error) {
	err := validateGenesisAccounts(cfg.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	chooserAddress, _, err := deployer.DeployProofOfBurn(blockchain)
	if err != nil {
		return nil, err
	}

	accountRegistryAddress, accountRegistryDeploymentBlock, accountRegistry, err := deployer.DeployAccountRegistry(blockchain, chooserAddress)
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

	registeredAccounts, err := RegisterGenesisAccounts(accountManager, cfg.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	populatedAccounts := AssignStateIDs(registeredAccounts)

	err = PopulateGenesisAccounts(storage, populatedAccounts)
	if err != nil {
		return nil, err
	}

	stateRoot, err := storage.StateTree.Root()
	if err != nil {
		return nil, err
	}

	contracts, err := rollup.DeployConfiguredRollup(blockchain, rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: stateRoot,
			BlocksToFinalise: models.NewUint256(uint64(cfg.BlocksToFinalise)),
		},
		Dependencies: rollup.Dependencies{
			AccountRegistry: accountRegistryAddress,
			Chooser:         chooserAddress,
		},
	})
	if err != nil {
		return nil, err
	}

	chainState := &models.ChainState{
		ChainID:                        blockchain.GetChainID(),
		AccountRegistry:                *accountRegistryAddress,
		AccountRegistryDeploymentBlock: *accountRegistryDeploymentBlock,
		TokenRegistry:                  contracts.TokenRegistryAddress,
		SpokeRegistry:                  contracts.SpokeRegistryAddress,
		DepositManager:                 contracts.DepositManagerAddress,
		WithdrawManager:                contracts.WithdrawManagerAddress,
		Rollup:                         contracts.RollupAddress,
		GenesisAccounts:                populatedAccounts,
		SyncedBlock:                    getInitialSyncedBlock(*accountRegistryDeploymentBlock),
	}

	return chainState, nil
}

func validateGenesisAccounts(accounts []models.GenesisAccount) error {
	for i := range accounts {
		if accounts[i].PublicKey == nil {
			return ErrNoPublicKeysInGenesisAccounts
		}
	}

	return nil
}
