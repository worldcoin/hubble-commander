package eth

import (
	"math/big"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/eth/deployer"
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
	events, unsubscribe, err := s.client.WatchTokenRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	err = s.client.RequestRegisterToken(s.client.ExampleTokenAddress)
	s.NoError(err)

	err = s.client.FinalizeRegisterToken(s.client.ExampleTokenAddress)
	s.NoError(err)

	var tokenID *big.Int
Outer:
	for {
		select {
		case event, ok := <-events:
			if !ok {
				s.Fail("Token registry event watcher is closed")
			}
			if event.TokenContract == s.client.ExampleTokenAddress {
				tokenID = event.TokenID
				break Outer
			}
		case <-time.After(deployer.ChainTimeout):
			s.Fail("Token registry event watcher timed out")
		}
	}
	s.Equal(uint64(0), tokenID.Uint64())
}

func TestRegisterTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTokenTestSuite))
}
