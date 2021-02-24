package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
)

type FrontendContracts struct {
	FrontendGeneric         *generic.FrontendGeneric
	FrontendTransfer        *transfer.FrontendTransfer
	FrontendMassMigration   *massmigration.FrontendMassMigration
	FrontendCreate2Transfer *create2transfer.FrontendCreate2Transfer
}

func DeployFrontend(sim *simulator.Simulator) (*FrontendContracts, error) {
	deployer := sim.Account

	_, _, genericContract, err := generic.DeployFrontendGeneric(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	_, _, transferContract, err := transfer.DeployFrontendTransfer(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	_, _, migrationContract, err := massmigration.DeployFrontendMassMigration(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	_, _, createContract, err := create2transfer.DeployFrontendCreate2Transfer(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	sim.Backend.Commit()

	return &FrontendContracts{
		FrontendGeneric:         genericContract,
		FrontendTransfer:        transferContract,
		FrontendMassMigration:   migrationContract,
		FrontendCreate2Transfer: createContract,
	}, nil
}
