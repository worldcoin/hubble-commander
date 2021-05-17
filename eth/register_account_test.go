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

func (s *RegisterAccountTestSuite) TestRegisterAccount() {
	events, unsubscribe, err := s.client.WatchRegistrations(nil)
	s.NoError(err)
	defer unsubscribe()

	publicKey := models.PublicKey{1, 2, 3}
	pubKeyID, err := s.client.RegisterAccount(&publicKey, events)
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)

	pubKeyID, err = s.client.RegisterAccount(&publicKey, events)
	s.NoError(err)
	s.Equal(uint32(1), *pubKeyID)
}

func TestRegisterAccountTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterAccountTestSuite))
}
