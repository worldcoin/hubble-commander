package encoder

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DecoderTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	testTypes *types.TestTypes
}

func (s *DecoderTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DecoderTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	test, err := deployer.DeployTest(sim)
	s.NoError(err)

	s.testTypes = test.TestTypes
}

func (s *DecoderTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *DecoderTestSuite) TestDecodeMeta() {
	input, err := s.testTypes.EncodeMeta(
		nil,
		big.NewInt(1),
		big.NewInt(2),
		utils.RandomAddress(),
		big.NewInt(30_000_000),
	)

	output, err := s.testTypes.DecodeMeta(nil, input)
	s.NoError(err)

	expectedMeta := models.Meta{
		BatchType:  uint8(output.BatchType.Uint64()),
		Size:       uint8(output.Size.Uint64()),
		Committer:  output.Committer,
		FinaliseOn: uint32(output.FinaliseOn.Uint64()),
	}

	meta := DecodeMeta(input)

	s.Equal(expectedMeta, meta)
}

func TestDecoderTestSuite(t *testing.T) {
	suite.Run(t, new(DecoderTestSuite))
}
