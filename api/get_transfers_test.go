package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.Storage
	db      *db.TestDB
	tree    *st.StateTree
}

func (s *GetTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransfersTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage, nil}
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)
}

func (s *GetTransfersTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransfersTestSuite) TestApi_GetTransfer() {
	account := models.Account{
		AccountIndex: 1,
		PublicKey:    models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	userStates := []models.UserState{
		{
			AccountIndex: account.AccountIndex,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
		{
			AccountIndex: 2,
			TokenIndex:   models.MakeUint256(2),
			Balance:      models.MakeUint256(500),
			Nonce:        models.MakeUint256(0),
		},
	}

	err = s.tree.Set(0, &userStates[0])
	s.NoError(err)
	err = s.tree.Set(1, &userStates[1])
	s.NoError(err)

	transfer := dto.Transfer{
		FromStateID: ref.Uint32(0),
		ToStateID:   ref.Uint32(1),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   []byte{1, 2, 3, 4},
	}

	hash, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)

	res, err := s.api.GetTransfer(*hash)
	s.NoError(err)

	s.Equal(models.Pending, res.Status)
}

// nolint:funlen
func (s *GetTransfersTestSuite) TestApi_GetTransfers() {
	account := models.Account{
		AccountIndex: 1,
		PublicKey:    models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	userStates := []models.UserState{
		{
			AccountIndex: account.AccountIndex,
			TokenIndex:   models.MakeUint256(1),
			Balance:      models.MakeUint256(420),
			Nonce:        models.MakeUint256(0),
		},
		{
			AccountIndex: 2,
			TokenIndex:   models.MakeUint256(2),
			Balance:      models.MakeUint256(500),
			Nonce:        models.MakeUint256(0),
		},
		{
			AccountIndex: account.AccountIndex,
			TokenIndex:   models.MakeUint256(25),
			Balance:      models.MakeUint256(1),
			Nonce:        models.MakeUint256(73),
		},
	}

	err = s.tree.Set(0, &userStates[0])
	s.NoError(err)
	err = s.tree.Set(1, &userStates[1])
	s.NoError(err)
	err = s.tree.Set(2, &userStates[2])
	s.NoError(err)

	transfers := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(1234)),
				FromStateID:          0,
				Amount:               models.MakeUint256(1),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(0),
				Signature:            []byte{1, 2, 3, 4, 5},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(2345)),
				FromStateID:          0,
				Amount:               models.MakeUint256(2),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(1),
				Signature:            []byte{2, 3, 4, 5, 6},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(3456)),
				FromStateID:          1,
				Amount:               models.MakeUint256(3),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(0),
				Signature:            []byte{3, 4, 5, 6, 7},
				IncludedInCommitment: nil,
			},
			ToStateID: 0,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(4567)),
				FromStateID:          0,
				Amount:               models.MakeUint256(2),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(2),
				Signature:            []byte{2, 3, 4, 5, 6},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
	}

	err = s.storage.AddTransfer(&transfers[0])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[1])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[2])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[3])
	s.NoError(err)

	userTransfers, err := s.api.GetTransfers(&account.PublicKey)
	s.NoError(err)

	s.Len(userTransfers, 3)
	s.Equal(userTransfers[0].Transfer.Hash, transfers[0].Hash)
	s.Equal(userTransfers[1].Transfer.Hash, transfers[1].Hash)
	s.Equal(userTransfers[2].Transfer.Hash, transfers[3].Hash)
}

func TestGetTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransfersTestSuite))
}
