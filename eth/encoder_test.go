package eth

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/frontendtransfer"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim      *testutils.Simulator
	contract *frontendtransfer.FrontendTransfer
}

func (s *EncoderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *EncoderTestSuite) SetupTest() {
	sim, err := testutils.NewSimulator()
	s.NoError(err)
	s.sim = sim

	_, _, contract, err := frontendtransfer.DeployFrontendTransfer(sim.Account, sim.Backend)
	s.NoError(err)
	sim.Backend.Commit()
	s.contract = contract
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeTransferZero() {
	tx := frontendtransfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.contract.Encode(&bind.CallOpts{Pending: false}, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeTransferNonZero() {
	tx := frontendtransfer.OffchainTransfer{
		TxType:    big.NewInt(1),
		FromIndex: big.NewInt(2),
		ToIndex:   big.NewInt(3),
		Amount:    big.NewInt(4),
		Fee:       big.NewInt(5),
		Nonce:     big.NewInt(6),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.contract.Encode(&bind.CallOpts{Pending: false}, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
