package simulator

import (
	"math/big"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SimulatorTestSuite struct {
	*require.Assertions
	suite.Suite
	sim *Simulator
}

func (s *SimulatorTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SimulatorTestSuite) SetupTest() {
	sim, err := NewSimulator()
	s.NoError(err)
	s.sim = sim
}

func (s *SimulatorTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *SimulatorTestSuite) TestNewSimulator() {
	_, _, contract, err := transfer.DeployFrontendTransfer(s.sim.Account, s.sim.Backend)
	s.NoError(err)

	s.sim.Backend.Commit()

	_, err = contract.Encode(nil, transfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	})
	s.NoError(err)
}

func (s *SimulatorTestSuite) TestNewAutominingSimulator() {
	sim, err := NewAutominingSimulator()
	s.NoError(err)

	s.True(sim.IsAutomineEnabled())
	sim.StopAutomine()
}

func (s *SimulatorTestSuite) TestStartAutomine() {
	stop := s.sim.StartAutomine()
	defer stop()
	time.Sleep(220 * time.Millisecond)
	s.Equal(uint64(2), s.sim.Backend.Blockchain().CurrentHeader().Number.Uint64())
}

func (s *SimulatorTestSuite) TestStopAutomine() {
	s.sim.StartAutomine()
	time.Sleep(120 * time.Millisecond)
	s.sim.StopAutomine()
	time.Sleep(100 * time.Millisecond)
	s.Equal(uint64(1), s.sim.Backend.Blockchain().CurrentHeader().Number.Uint64())
}

func (s *SimulatorTestSuite) TestClose_StopsAutomine() {
	s.sim.StartAutomine()
	time.Sleep(120 * time.Millisecond)
	s.sim.Close()
	time.Sleep(100 * time.Millisecond)
	s.Equal(uint64(1), s.sim.Backend.Blockchain().CurrentHeader().Number.Uint64())
}

func (s *SimulatorTestSuite) TestGetLatestBlockNumber() {
	blockNumberBefore, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	s.sim.Backend.Commit()

	blockNumberAfter, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	expectedBlockNumber := *blockNumberBefore + 1
	s.Equal(expectedBlockNumber, *blockNumberAfter)
}

func TestSimulatorTestSuite(t *testing.T) {
	suite.Run(t, new(SimulatorTestSuite))
}
