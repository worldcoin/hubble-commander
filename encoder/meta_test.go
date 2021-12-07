package encoder

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/test/types"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MetaTestSuite struct {
	*require.Assertions
	suite.Suite
	sim       *simulator.Simulator
	testTypes *types.TestTypes
}

func (s *MetaTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MetaTestSuite) SetupTest() {
	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	test, err := testutils.DeployTest(sim)
	s.NoError(err)

	s.testTypes = test.TestTypes
}

func (s *MetaTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *MetaTestSuite) TestDecodeMeta() {
	input, err := s.testTypes.EncodeMeta(
		nil,
		big.NewInt(1),
		big.NewInt(2),
		utils.RandomAddress(),
		big.NewInt(30_000_000),
	)
	s.NoError(err)

	output, err := s.testTypes.DecodeMeta(nil, input)
	s.NoError(err)

	expectedMeta := models.BatchMeta{
		BatchType:  batchtype.BatchType(output.BatchType.Uint64()),
		Size:       uint8(output.Size.Uint64()),
		Committer:  output.Committer,
		FinaliseOn: uint32(output.FinaliseOn.Uint64()),
	}

	meta := DecodeMeta(input)

	s.Equal(expectedMeta, meta)
}

func TestMetaTestSuite(t *testing.T) {
	suite.Run(t, new(MetaTestSuite))
}
