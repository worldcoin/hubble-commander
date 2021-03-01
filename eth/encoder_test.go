package eth

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim      *simulator.Simulator
	transfer *transfer.FrontendTransfer
	generic  *generic.FrontendGeneric
}

func (s *EncoderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *EncoderTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	s.transfer = contracts.FrontendTransfer
	s.generic = contracts.FrontendGeneric
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeTransferZero() {
	tx := transfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.transfer.Encode(&bind.CallOpts{Pending: false}, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeTransferNonZero() {
	tx := transfer.OffchainTransfer{
		TxType:    big.NewInt(1),
		FromIndex: big.NewInt(2),
		ToIndex:   big.NewInt(3),
		Amount:    big.NewInt(4),
		Fee:       big.NewInt(5),
		Nonce:     big.NewInt(6),
	}
	bytes, err := EncodeTransfer(tx)
	s.NoError(err)
	expected, err := s.transfer.Encode(&bind.CallOpts{Pending: false}, tx)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeUserState() {
	state := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(2),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	expectedState := generic.TypesUserState{
		PubkeyID: big.NewInt(1),
		TokenID:  big.NewInt(2),
		Balance:  big.NewInt(420),
		Nonce:    big.NewInt(0),
	}
	bytes, err := EncodeUserState(state)
	s.NoError(err)
	expected, err := s.generic.Encode(&bind.CallOpts{Pending: false}, expectedState)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
