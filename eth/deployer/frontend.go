package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/Worldcoin/hubble-commander/eth/chain"
)

type FrontendContracts struct {
	FrontendGeneric         *generic.FrontendGeneric
	FrontendTransfer        *transfer.FrontendTransfer
	FrontendMassMigration   *massmigration.FrontendMassMigration
	FrontendCreate2Transfer *create2transfer.FrontendCreate2Transfer
}

func DeployFrontend(c chain.Connection) (*FrontendContracts, error) {
	_, _, genericContract, err := generic.DeployFrontendGeneric(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, transferContract, err := transfer.DeployFrontendTransfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, migrationContract, err := massmigration.DeployFrontendMassMigration(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, createContract, err := create2transfer.DeployFrontendCreate2Transfer(c.GetAccount(), c.GetBackend())
	if err != nil {
		return nil, err
	}

	c.GetBackend().Commit()

	return &FrontendContracts{
		FrontendGeneric:         genericContract,
		FrontendTransfer:        transferContract,
		FrontendMassMigration:   migrationContract,
		FrontendCreate2Transfer: createContract,
	}, nil
}
