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

type GetUserStateProofTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
}

func (s *GetUserStateProofTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStateProofTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{
		storage: s.storage.Storage,
		cfg:     &config.APIConfig{EnableProofMethods: true},
	}
}

func (s *GetUserStateProofTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetUserStateProofTestSuite) TestGetUserState() {
	leaf := models.StateLeaf{
		StateID: 0,
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(420),
			Nonce:    models.MakeUint256(0),
		},
	}
	_, err := s.api.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	witness, err := s.api.storage.StateTree.GetLeafWitness(leaf.StateID)
	s.NoError(err)

	dtoUserState := dto.MakeUserState(&leaf.UserState)

	expectedUserStateProof := &dto.StateMerkleProof{
		UserState: &dtoUserState,
		Witness:   witness,
	}
	userStateProof, err := s.api.GetUserStateProof(leaf.StateID)
	s.NoError(err)
	s.Equal(expectedUserStateProof, userStateProof)
}

func (s *GetUserStateProofTestSuite) TestGetUserState_NonexistentStateLeaf() {
	_, err := s.api.GetUserStateProof(1)
	s.Equal(&APIError{
		Code:    50003,
		Message: "user state inclusion proof could not be generated",
	}, err)
}

func TestGetUserStateProofTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStateProofTestSuite))
}
