package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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
	tokenAddress := s.deployToken()

	tokenID, err := s.client.RegisterTokenAndWait(tokenAddress)
	s.NoError(err)
	s.Equal(tokenID, models.NewUint256(1))
}

func (s *RegisterTokenTestSuite) TestGetRegisteredToken_ReturnsCorrectToken() {
	registeredToken, err := s.client.GetRegisteredToken(models.NewUint256(0))
	s.NoError(err)
	s.Equal(registeredToken.Contract, s.client.ExampleTokenAddress)
}

func (s *RegisterTokenTestSuite) deployToken() common.Address {
	tokenAddress, tokenTx, _, err := customtoken.DeployTestCustomToken(
		s.client.GetAccount(),
		s.client.GetBackend(),
		"NewToken",
		"NEW",
	)
	s.NoError(err)
	_, err = chain.WaitToBeMined(s.client.Backend, tokenTx)
	s.NoError(err)

	return tokenAddress
}

func TestRegisterTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTokenTestSuite))
}
