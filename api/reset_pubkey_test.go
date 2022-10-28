// we are testing `api/admin/reset_pubkey` but we call `api` methods so
// the test lives here to break the import loop
package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/api/admin"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ResetPubkeyTestSuite struct {
	*require.Assertions
	suite.Suite
	hubbleAPI *API
	adminAPI  *admin.API
	client    *eth.TestClient
	storage   *storage.TestStorage
	tree      *storage.AccountTree
}

func (s *ResetPubkeyTestSuite) makePubKey(x int64) models.PublicKey {
	ints := [4]*big.Int{
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(x),
	}

	return models.MakePublicKeyFromInts(ints)
}

// this seems like a lovely place for a prop test

func (s *ResetPubkeyTestSuite) Test_RoundTripToEmpty() {
	zeroRoot := merkletree.GetZeroHash(32)
	pubKey := s.makePubKey(1)

	// the tree had better start with our known zero root

	oldRootPtr, err := s.tree.Root()
	s.NoError(err)
	s.Equal(zeroRoot, *oldRootPtr)

	// once we add a key the root will change

	oneResult, err := s.adminAPI.ResetPubkey(contextWithAuthKey(), 0, &pubKey)
	s.NoError(err)
	s.NotNil(oneResult)
	s.NotEqual(oldRootPtr, oneResult.NewAccountTreeRoot)
	s.Nil(oneResult.OldPubKey) // the tree used to be empty, there was no old key

	// now set it back to empty and check that the new root equals oldRootPtr

	emptyResult, err := s.adminAPI.ResetPubkey(contextWithAuthKey(), 0, nil)
	s.NoError(err)
	s.Equal(*oldRootPtr, emptyResult.NewAccountTreeRoot)
	s.Equal(&pubKey, emptyResult.OldPubKey)
}

func (s *ResetPubkeyTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ResetPubkeyTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.adminAPI = admin.NewTestAPI(
		&config.APIConfig{AuthenticationKey: authKeyValue},
		s.storage.Storage,
		s.client.Client,
	)
	s.hubbleAPI = NewTestAPI(
		s.storage.Storage,
		s.client.Client,
	)

	s.tree = s.storage.AccountTree
}

func (s *ResetPubkeyTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func TestResetPubkeyTestSuite(t *testing.T) {
	suite.Run(t, new(ResetPubkeyTestSuite))
}
