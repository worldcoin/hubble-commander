package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func DeployAccountRegistry(c ChainConnection) (*common.Address, *accountregistry.AccountRegistry, error) {
	accountRegistryAddress, _, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	c.Commit()

	return &accountRegistryAddress, accountRegistry, nil
}
