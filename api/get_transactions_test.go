package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionsTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.Storage
	db      *db.TestDB
	tree    *st.StateTree
}

func (s *GetTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage}
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)
}

func (s *GetTransactionsTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransactionsTestSuite) XDDD() {
	incomingTx := models.IncomingTransaction{
		FromIndex: models.NewUint256(1),
		ToIndex:   models.NewUint256(2),
		Amount:    models.NewUint256(50),
		Fee:       models.NewUint256(10),
		Nonce:     models.NewUint256(0),
		Signature: []byte{1, 2, 3, 4},
	}

	hash, err := s.api.SendTransaction(incomingTx)
	s.NoError(err)

	res, err := s.api.GetTransaction(*hash)
	s.NoError(err)

	s.Equal(models.Pending, res.Status)
}

func (s *GetTransactionsTestSuite) TestApi_GetTransactions() {
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

	transactions := []models.Transaction{
		{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			FromIndex:            models.MakeUint256(0),
			ToIndex:              models.MakeUint256(1),
			Amount:               models.MakeUint256(1),
			Fee:                  models.MakeUint256(5),
			Nonce:                models.MakeUint256(0),
			Signature:            []byte{1, 2, 3, 4, 5},
			IncludedInCommitment: nil,
		},
		{
			Hash:                 common.BigToHash(big.NewInt(2345)),
			FromIndex:            models.MakeUint256(0),
			ToIndex:              models.MakeUint256(1),
			Amount:               models.MakeUint256(2),
			Fee:                  models.MakeUint256(5),
			Nonce:                models.MakeUint256(1),
			Signature:            []byte{2, 3, 4, 5, 6},
			IncludedInCommitment: nil,
		},
		{
			Hash:                 common.BigToHash(big.NewInt(3456)),
			FromIndex:            models.MakeUint256(1),
			ToIndex:              models.MakeUint256(0),
			Amount:               models.MakeUint256(3),
			Fee:                  models.MakeUint256(5),
			Nonce:                models.MakeUint256(0),
			Signature:            []byte{3, 4, 5, 6, 7},
			IncludedInCommitment: nil,
		},
		{
			Hash:                 common.BigToHash(big.NewInt(4567)),
			FromIndex:            models.MakeUint256(0),
			ToIndex:              models.MakeUint256(1),
			Amount:               models.MakeUint256(2),
			Fee:                  models.MakeUint256(5),
			Nonce:                models.MakeUint256(2),
			Signature:            []byte{2, 3, 4, 5, 6},
			IncludedInCommitment: nil,
		},
	}

	err = s.storage.AddTransaction(&transactions[0])
	s.NoError(err)
	err = s.storage.AddTransaction(&transactions[1])
	s.NoError(err)
	err = s.storage.AddTransaction(&transactions[2])
	s.NoError(err)
	err = s.storage.AddTransaction(&transactions[3])
	s.NoError(err)

	userTransactions, err := s.api.GetTransactions(&account.PublicKey)
	s.NoError(err)

	s.Len(userTransactions, 3)
	s.Equal(userTransactions[0].Transaction.Hash, transactions[0].Hash)
	s.Equal(userTransactions[1].Transaction.Hash, transactions[1].Hash)
	s.Equal(userTransactions[2].Transaction.Hash, transactions[3].Hash)
}

func TestGetTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionsTestSuite))
}
