package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/ethereum/go-ethereum/common"
)

func DeployAccountRegistry(d Deployer) (*common.Address, *accountregistry.AccountRegistry, error) {
	accountRegistryAddress, _, accountRegistry, err := accountregistry.DeployAccountRegistry(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, nil, err
	}

	d.Commit()

	return &accountRegistryAddress, accountRegistry, err
}
