package eth

import (
	"log"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	requestEvents, requestUnsubscribe, err := s.client.WatchTokenRegistrationRequests(&bind.WatchOpts{Start: nil})
	s.NoError(err)
	defer requestUnsubscribe()
	events, unsubscribe, err := s.client.WatchTokenRegistrations(&bind.WatchOpts{Start: nil})
	s.NoError(err)
	defer unsubscribe()

	tokenContract := utils.RandomAddress()
	log.Printf("Token contract %v\n", tokenContract)

	s.client.RequestRegisterToken(tokenContract, requestEvents)
	tokenID, err := s.client.FinalizeRegisterToken(tokenContract, events)
	s.NoError(err)
	s.Equal(uint32(2), *tokenID)
}

func TestRegisterTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTokenTestSuite))
}
