package encoder

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EncoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	generic   *generic.FrontendGeneric
	testTx    *testtx.TestTx
	testTypes *types.TestTypes
}

func (s *EncoderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *EncoderTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	frontend, err := deployer.DeployFrontend(sim)
	s.NoError(err)
	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.generic = frontend.FrontendGeneric
	s.testTx = test.TestTx
	s.testTypes = test.TestTypes
}

func (s *EncoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *EncoderTestSuite) TestEncodeUserState() {
	state := generic.TypesUserState{
		PubkeyID: big.NewInt(1),
		TokenID:  big.NewInt(2),
		Balance:  big.NewInt(420),
		Nonce:    big.NewInt(0),
	}
	bytes, err := EncodeUserState(state)
	s.NoError(err)

	expected, err := s.generic.Encode(nil, state)
	s.NoError(err)
	s.Equal(expected, bytes)
}

func (s *EncoderTestSuite) TestEncodeDecimal() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	expected, err := s.testTx.TestEncodeDecimal(nil, &num.Int)
	s.NoError(err)

	s.Equal(uint16(expected.Uint64()), encoded)
}

func (s *EncoderTestSuite) TestEncodeAndDecodeDecimal() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	decoded := DecodeDecimal(encoded)

	s.Equal(num, decoded)
}

func (s *EncoderTestSuite) TestCommitmentBodyHash() {
	accountRoot := utils.RandomHash()
	signature := models.MakeRandomSignature()
	feeReceiver := models.MakeUint256(1234)
	txs := utils.RandomBytes(32)

	expectedHash, err := s.testTypes.HashTransferBody(nil, types.TypesTransferBody{
		AccountRoot: accountRoot,
		Signature:   signature.BigInts(),
		FeeReceiver: &feeReceiver.Int,
		Txs:         txs,
	})
	s.NoError(err)

	commitment := models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      txs,
		FeeReceiver:       uint32(feeReceiver.Uint64()),
		CombinedSignature: signature,
		AccountTreeRoot:   &accountRoot,
	}

	s.Equal(expectedHash[:], commitment.BodyHash().Bytes())
}

func TestEncoderTestSuite(t *testing.T) {
	suite.Run(t, new(EncoderTestSuite))
}
