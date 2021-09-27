package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisterAccountTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *RegisterAccountTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisterAccountTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *RegisterAccountTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *RegisterAccountTestSuite) TestRegisterAccountAndWait() {
	publicKey := models.PublicKey{1, 2, 3}
	pubKeyID, err := s.client.RegisterAccountAndWait(&publicKey)
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)

	pubKeyID, err = s.client.RegisterAccountAndWait(&publicKey)
	s.NoError(err)
	s.Equal(uint32(1), *pubKeyID)
}

func (s *RegisterAccountTestSuite) TestGetNextSingleRegistrationPubKeyID() {
	events, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	lastPubKeyID := new(uint32)
	for i := 0; i < 3; i++ {
		lastPubKeyID, err = s.client.RegisterAccount(&models.PublicKey{1, 2, 3}, events)
		s.NoError(err)
	}

	pubKeyID, err := s.client.GetNextSingleRegistrationPubKeyID()
	s.NoError(err)
	s.EqualValues(*lastPubKeyID+1, *pubKeyID)
}

func (s *RegisterAccountTestSuite) TestGetNextSingleRegistrationPubKeyID_NoAccounts() {
	pubKeyID, err := s.client.GetNextSingleRegistrationPubKeyID()
	s.NoError(err)
	s.EqualValues(0, *pubKeyID)
}

func TestRegisterAccountTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterAccountTestSuite))
}
