package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserStateProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	teardown func() error
}

func (s *GetUserStateProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStateProofTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.api = &API{storage: testStorage.Storage}
}

func (s *GetUserStateProofTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *GetUserStateProofTestSuite) TestGetUserStates() {
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

	witness, err := s.api.storage.StateTree.GetLeafWitness(0)
	s.NoError(err)

	expectedUserStateProof := &dto.StateMerkleProof{
		StateMerkleProof: models.StateMerkleProof{
			UserState: &leaf.UserState,
			Witness:   witness,
		},
	}
	userStateProof, err := s.api.GetUserStateProof(0)

	s.NoError(err)
	s.Equal(expectedUserStateProof, userStateProof)
}

func TestGetUserStateProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStateProofTestSuite))
}
