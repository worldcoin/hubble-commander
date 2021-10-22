package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetPublicKeyProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	teardown func() error
}

func (s *GetPublicKeyProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetPublicKeyProofTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.api = &API{storage: testStorage.Storage}
}

func (s *GetPublicKeyProofTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *GetPublicKeyProofTestSuite) TestGetPublicKeyProofByPubKeyID() {
	account := models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.api.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	leaf := models.StateLeaf{
		StateID: 0,
		UserState: models.UserState{
			PubKeyID: account.PubKeyID,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
	}
	_, err = s.api.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	witness, err := s.api.storage.AccountTree.GetWitness(account.PubKeyID)
	s.NoError(err)

	expectedPublicKeyProof := &dto.PublicKeyProof{
		PublicKeyProof: models.PublicKeyProof{
			PublicKey: &account.PublicKey,
			Witness:   witness,
		},
	}
	publicKeyProof, err := s.api.GetPublicKeyProofByPubKeyID(account.PubKeyID)
	s.NoError(err)
	s.Equal(expectedPublicKeyProof, publicKeyProof)
}

func (s *GetPublicKeyProofTestSuite) TestGetPublicKeyProofByPubKeyID_NonexistentAccount() {
	_, err := s.api.GetPublicKeyProofByPubKeyID(1)
	s.Equal(&APIError{
		Code:    99005,
		Message: "public key proof not found",
	}, err)
}

func TestGetPublicKeyProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetPublicKeyProofTestSuite))
}
