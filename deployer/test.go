package deployer

import (
	"github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
)

type TestContracts struct {
	TestTx    *tx.TestTx
	TestTypes *types.TestTypes
}

func DeployTest(sim *simulator.Simulator) (*TestContracts, error) {
	deployer := sim.Account

	_, _, testTx, err := tx.DeployTestTx(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	_, _, testTypes, err := types.DeployTestTypes(deployer, sim.Backend)
	if err != nil {
		return nil, err
	}

	sim.Backend.Commit()

	return &TestContracts{
		TestTx:    testTx,
		TestTypes: testTypes,
	}, nil
}
