package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserStatesTestSuite struct {
	*require.Assertions
	suite.Suite
	api *API
	db  *db.TestDB
	tx  *models.Transaction
}

func (s *GetUserStatesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStatesTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{nil, storage}
	s.db = testDB

	tx := &models.Transaction{
		FromIndex: *models.NewUint256(1),
		ToIndex:   *models.NewUint256(2),
		Amount:    *models.NewUint256(50),
		Fee:       *models.NewUint256(10),
		Nonce:     *models.NewUint256(0),
		Signature: []byte{1, 2, 3, 4},
	}

	s.tx = tx
}

func (s *GetUserStatesTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetUserStatesTestSuite) TestApi_GetTransaction() {
	account := models.Account{
		AccountIndex: 1,
		PublicKey:    models.PublicKey{1, 2, 3},
	}

	err := s.api.storage.AddAccount(&account)
	s.NoError(err)

	leafs := []models.StateLeaf{
		{
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				AccountIndex: account.AccountIndex,
				TokenIndex:   models.MakeUint256(1),
				Balance:      models.MakeUint256(420),
				Nonce:        models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				AccountIndex: account.AccountIndex,
				TokenIndex:   models.MakeUint256(2),
				Balance:      models.MakeUint256(500),
				Nonce:        models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{3, 4, 5, 6, 7}),
			UserState: models.UserState{
				AccountIndex: account.AccountIndex,
				TokenIndex:   models.MakeUint256(25),
				Balance:      models.MakeUint256(1),
				Nonce:        models.MakeUint256(73),
			},
		},
	}

	err = s.api.storage.AddStateLeaf(&leafs[0])
	s.NoError(err)
	err = s.api.storage.AddStateLeaf(&leafs[1])
	s.NoError(err)
	err = s.api.storage.AddStateLeaf(&leafs[2])
	s.NoError(err)

	userStates, err := s.api.GetUserStates(&account.PublicKey)
	s.NoError(err)

	s.Len(userStates, 3)
	s.Equal(leafs[0].UserState, userStates[0])
	s.Equal(leafs[1].UserState, userStates[1])
	s.Equal(leafs[2].UserState, userStates[2])
}

func TestGetUserStatesTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStatesTestSuite))
}
