package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisterTokenTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *RegisterTokenTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisterTokenTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *RegisterTokenTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *RegisterTokenTestSuite) TestRegisterTokenAndWait_ReturnsCorrectToken() {
	tokenID, err := s.client.RegisterTokenAndWait(s.client.ExampleTokenAddress)
	s.NoError(err)
	s.Equal(tokenID, models.NewUint256(0))
}

func TestRegisterTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTokenTestSuite))
}
