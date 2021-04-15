package deployer

import (
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func DeployAccountRegistry(c ChainConnection) (*common.Address, *accountregistry.AccountRegistry, error) {
	log.Println("Deploying AccountRegistry")
	accountRegistryAddress, tx, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}
	c.Commit()

	_, err = WaitToBeMined(c.GetBackend(), tx)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return &accountRegistryAddress, accountRegistry, nil
}
