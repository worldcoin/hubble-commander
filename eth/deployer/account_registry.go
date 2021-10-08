package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func DeployAccountRegistry(c chain.Connection, chooser *common.Address) (
	*common.Address, *uint64, *accountregistry.AccountRegistry, error,
) {
	log.Println("Deploying AccountRegistry")
	accountRegistryAddress, tx, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend(), *chooser)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	c.Commit()
	txReceipt, err := chain.WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	deploymentBlockNumber := txReceipt.BlockNumber.Uint64()

	return &accountRegistryAddress, &deploymentBlockNumber, accountRegistry, nil
}
