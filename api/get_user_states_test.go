package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserStatesTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	teardown func() error
}

func (s *GetUserStatesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStatesTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.api = &API{storage: testStorage.Storage}
}

func (s *GetUserStatesTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *GetUserStatesTestSuite) TestGetUserStates() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err := s.api.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	leaves := []models.StateLeaf{
		{
			StateID:  0,
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				PubKeyID: accounts[0].PubKeyID,
				TokenID:  models.MakeUint256(1),
				Balance:  models.MakeUint256(420),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID:  1,
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				PubKeyID: accounts[1].PubKeyID,
				TokenID:  models.MakeUint256(2),
				Balance:  models.MakeUint256(500),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID:  2,
			DataHash: common.BytesToHash([]byte{3, 4, 5, 6, 7}),
			UserState: models.UserState{
				PubKeyID: accounts[0].PubKeyID,
				TokenID:  models.MakeUint256(25),
				Balance:  models.MakeUint256(1),
				Nonce:    models.MakeUint256(73),
			},
		},
	}
	for i := range leaves {
		err := s.api.storage.UpsertStateLeaf(&leaves[i])
		s.NoError(err)
	}

	userStates, err := s.api.GetUserStates(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(userStates, 3)
	s.Contains(userStates, dto.UserState{
		StateID:   0,
		UserState: leaves[0].UserState,
	})
	s.Contains(userStates, dto.UserState{
		StateID:   1,
		UserState: leaves[1].UserState,
	})
	s.Contains(userStates, dto.UserState{
		StateID:   2,
		UserState: leaves[2].UserState,
	})
}

func TestGetUserStatesTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStatesTestSuite))
}
