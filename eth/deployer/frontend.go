package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/create2transfer"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/massmigration"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
)

type FrontendContracts struct {
	FrontendGeneric         *generic.FrontendGeneric
	FrontendTransfer        *transfer.FrontendTransfer
	FrontendMassMigration   *massmigration.FrontendMassMigration
	FrontendCreate2Transfer *create2transfer.FrontendCreate2Transfer
}

func DeployFrontend(d ChainConnection) (*FrontendContracts, error) {
	_, _, genericContract, err := generic.DeployFrontendGeneric(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, transferContract, err := transfer.DeployFrontendTransfer(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, migrationContract, err := massmigration.DeployFrontendMassMigration(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	_, _, createContract, err := create2transfer.DeployFrontendCreate2Transfer(d.TransactionOpts(), d.GetBackend())
	if err != nil {
		return nil, err
	}

	d.Commit()

	return &FrontendContracts{
		FrontendGeneric:         genericContract,
		FrontendTransfer:        transferContract,
		FrontendMassMigration:   migrationContract,
		FrontendCreate2Transfer: createContract,
	}, nil
}
