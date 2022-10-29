package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBlocksToFinaliseTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *GetBlocksToFinaliseTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBlocksToFinaliseTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *GetBlocksToFinaliseTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *GetBlocksToFinaliseTestSuite) TestGetBlocksToFinalise() {
	expected := int64(5760)
	s.Equal(expected, int64(config.DefaultBlocksToFinalise))

	blocksToFinalise, err := s.client.GetBlocksToFinalise()
	s.NoError(err)
	s.Equal(expected, *blocksToFinalise)
	s.Equal(expected, *s.client.blocksToFinalise)
}

func TestGetBlocksToFinaliseTestSuite(t *testing.T) {
	suite.Run(t, new(GetBlocksToFinaliseTestSuite))
}
