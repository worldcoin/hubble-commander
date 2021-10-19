package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositSubtreeDepthTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *DepositSubtreeDepthTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositSubtreeDepthTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *DepositSubtreeDepthTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *DepositSubtreeDepthTestSuite) TestGetDepositSubtreeDepth() {
	expectedSubtreeDepth, err := s.client.DepositManager.ParamMaxSubtreeDepth(&bind.CallOpts{})
	s.NoError(err)

	subtreeDepth, err := s.client.GetDepositSubtreeDepth()
	s.NoError(err)
	s.Equal(uint8(expectedSubtreeDepth.Uint64()), *subtreeDepth)
	s.Equal(s.client.depositSubtreeDepth, subtreeDepth)
}

func TestDepositSubtreeDepthTestSuite(t *testing.T) {
	suite.Run(t, new(DepositSubtreeDepthTestSuite))
}
