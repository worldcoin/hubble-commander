package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/ethereum/go-ethereum/common"
)

func DeployAccountRegistry(c ChainConnection) (*common.Address, *accountregistry.AccountRegistry, error) {
	accountRegistryAddress, _, accountRegistry, err := accountregistry.DeployAccountRegistry(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, nil, err
	}

	c.Commit()

	return &accountRegistryAddress, accountRegistry, err
}
