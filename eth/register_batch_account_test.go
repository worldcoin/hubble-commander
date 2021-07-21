package eth

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisterBatchAccountTestSuite struct {
	*require.Assertions
	suite.Suite
	client *TestClient
}

func (s *RegisterBatchAccountTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisterBatchAccountTestSuite) SetupTest() {
	client, err := NewTestClient()
	s.NoError(err)
	s.client = client
}

func (s *RegisterBatchAccountTestSuite) TearDownTest() {
	s.client.Close()
}

func (s *RegisterBatchAccountTestSuite) TestRegisterBatchAccount() {
	events, unsubscribe, err := s.client.WatchBatchAccountRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	var publicKeys [16]models.PublicKey
	expectedPubKeyIDs := make([]uint32, 16)
	for i := range publicKeys {
		publicKeys[i] = models.PublicKey{1, 1, byte(i)}
		expectedPubKeyIDs[i] = uint32(batchAccountOffset + i)
	}

	pubKeyIDs, err := s.client.RegisterBatchAccount(publicKeys, events)
	s.NoError(err)
	s.Len(pubKeyIDs, 16)
	s.Equal(expectedPubKeyIDs, pubKeyIDs)

	rightIndex, err := s.client.AccountRegistry.LeafIndexRight(&bind.CallOpts{})
	s.NoError(err)
	s.EqualValues(16, rightIndex.Uint64())
}

func TestRegisterBatchAccountTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterBatchAccountTestSuite))
}
