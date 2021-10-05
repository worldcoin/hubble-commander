package commander

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	log "github.com/sirupsen/logrus"
)

var ErrNoPublicKeysInGenesisAccounts = fmt.Errorf("genesis accounts for deployment require public keys")

func Deploy(cfg *config.Config, chain deployer.ChainConnection) (chainSpec *string, err error) {
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
		"Bootstrapping genesis state with %d accounts on chainId = %s",
		len(cfg.Bootstrap.GenesisAccounts),
		cfg.Ethereum.ChainID,
	)
	chainState, err := deployContractsAndSetupGenesisState(tempStorage.Storage, chain, cfg.Bootstrap)
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
	chain deployer.ChainConnection,
	config *config.BootstrapConfig,
) (*models.ChainState, error) {
	err := validateGenesisAccounts(config.GenesisAccounts)
	if err != nil {
		return nil, err
	}

	proofOfBurnAddress, _, err := deployer.DeployProofOfBurn(chain)
	if err != nil {
		return nil, err
	}

	accountRegistryAddress, accountRegistryDeploymentBlock, accountRegistry, err := deployer.DeployAccountRegistry(chain, proofOfBurnAddress)
	if err != nil {
		return nil, err
	}

	registeredAccounts, err := RegisterGenesisAccounts(chain.GetAccount(), accountRegistry, config.GenesisAccounts)
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

	contracts, err := rollup.DeployConfiguredRollup(chain, rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: stateRoot,
			BlocksToFinalise: models.NewUint256(config.BlocksToFinalise),
		},
		Dependencies: rollup.Dependencies{AccountRegistry: accountRegistryAddress},
	})
	if err != nil {
		return nil, err
	}

	chainState := &models.ChainState{
		ChainID:         chain.GetChainID(),
		AccountRegistry: *accountRegistryAddress,
		DeploymentBlock: *accountRegistryDeploymentBlock,
		Rollup:          contracts.RollupAddress,
		GenesisAccounts: populatedAccounts,
		SyncedBlock:     getInitialSyncedBlock(*accountRegistryDeploymentBlock),
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
