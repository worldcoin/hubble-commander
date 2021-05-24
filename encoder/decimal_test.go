package encoder

import (
	"testing"

	testtx "github.com/Worldcoin/hubble-commander/contracts/test/tx"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DecimalTestSuite struct {
	*require.Assertions
	suite.Suite
	sim    *simulator.Simulator
	testTx *testtx.TestTx
}

func (s *DecimalTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DecimalTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.testTx = test.TestTx
}

func (s *DecimalTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *DecimalTestSuite) TestEncodeDecimal() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	expected, err := s.testTx.TestEncodeDecimal(nil, num.ToBig())
	s.NoError(err)

	s.Equal(uint16(expected.Uint64()), encoded)
}

func (s *DecimalTestSuite) TestEncodeAndDecodeDecimal() {
	num := models.MakeUint256(123400000)
	encoded, err := EncodeDecimal(num)
	s.NoError(err)

	decoded := DecodeDecimal(encoded)

	s.Equal(num, decoded)
}

func TestDecimalTestSuite(t *testing.T) {
	suite.Run(t, new(DecimalTestSuite))
}
