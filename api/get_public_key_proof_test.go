package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetPublicKeyProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
}

func (s *GetPublicKeyProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetPublicKeyProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		storage: s.storage.Storage,
		cfg:     &config.APIConfig{EnableProofMethods: true},
	}
}

func (s *GetPublicKeyProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
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
		PublicKey: &account.PublicKey,
		Witness:   witness,
	}
	publicKeyProof, err := s.api.GetPublicKeyProofByPubKeyID(account.PubKeyID)
	s.NoError(err)
	s.Equal(expectedPublicKeyProof, publicKeyProof)
}

func (s *GetPublicKeyProofTestSuite) TestGetPublicKeyProofByPubKeyID_NonexistentAccount() {
	_, err := s.api.GetPublicKeyProofByPubKeyID(1)
	s.Equal(&APIError{
		Code:    50002,
		Message: "public key inclusion proof could not be generated",
	}, err)
}

func TestGetPublicKeyProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetPublicKeyProofTestSuite))
}
