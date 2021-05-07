package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db/postgres"
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
	api *API
	db  *postgres.TestDB
}

func (s *GetUserStatesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStatesTestSuite) SetupTest() {
	testDB, err := postgres.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{nil, storage, nil}
	s.db = testDB
}

func (s *GetUserStatesTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetUserStatesTestSuite) TestGetUserStates() {
	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}

	err := s.api.storage.AddAccountIfNotExists(&accounts[0])
	s.NoError(err)
	err = s.api.storage.AddAccountIfNotExists(&accounts[1])
	s.NoError(err)

	leaves := []models.StateLeaf{
		{
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				PubKeyID:   accounts[0].PubKeyID,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(420),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				PubKeyID:   accounts[1].PubKeyID,
				TokenIndex: models.MakeUint256(2),
				Balance:    models.MakeUint256(500),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{3, 4, 5, 6, 7}),
			UserState: models.UserState{
				PubKeyID:   accounts[0].PubKeyID,
				TokenIndex: models.MakeUint256(25),
				Balance:    models.MakeUint256(1),
				Nonce:      models.MakeUint256(73),
			},
		},
	}

	err = s.api.storage.AddStateLeaf(&leaves[0])
	s.NoError(err)
	err = s.api.storage.AddStateLeaf(&leaves[1])
	s.NoError(err)
	err = s.api.storage.AddStateLeaf(&leaves[2])
	s.NoError(err)

	path, err := models.NewMerklePath("00")
	s.NoError(err)
	err = s.api.storage.UpsertStateNode(&models.StateNode{
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		MerklePath: *path,
	})
	s.NoError(err)

	path, err = models.NewMerklePath("01")
	s.NoError(err)
	err = s.api.storage.UpsertStateNode(&models.StateNode{
		DataHash:   common.BytesToHash([]byte{2, 3, 4, 5, 6}),
		MerklePath: *path,
	})
	s.NoError(err)

	path, err = models.NewMerklePath("10")
	s.NoError(err)
	err = s.api.storage.UpsertStateNode(&models.StateNode{
		DataHash:   common.BytesToHash([]byte{3, 4, 5, 6, 7}),
		MerklePath: *path,
	})
	s.NoError(err)

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
