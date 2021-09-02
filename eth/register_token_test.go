package eth

import (
	"testing"

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
	err := s.client.RequestRegisterToken(s.client.CustomTokenAddress)
	if err != nil {
	}

	tokenID, err := s.client.FinalizeRegisterToken(s.client.CustomTokenAddress)
	s.NoError(err)
	s.Equal(uint32(0), *tokenID)
}

func TestRegisterTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTokenTestSuite))
}
