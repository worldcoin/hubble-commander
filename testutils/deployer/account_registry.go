package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/common"
)

func DeployAccountRegistry(sim *simulator.Simulator) (*common.Address, *accountregistry.AccountRegistry, error) {
	accountRegistryAddress, _, accountRegistry, err := accountregistry.DeployAccountRegistry(sim.Account, sim.Backend)
	if err != nil {
		return nil, nil, err
	}

	sim.Backend.Commit()

	return &accountRegistryAddress, accountRegistry, err
}
