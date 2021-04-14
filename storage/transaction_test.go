package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	tx = models.Transaction{
		Hash:                 common.BigToHash(big.NewInt(1234)),
		FromIndex:            1,
		ToIndex:              2,
		Amount:               models.MakeUint256(1000),
		Fee:                  models.MakeUint256(100),
		Nonce:                models.MakeUint256(0),
		Signature:            []byte{1, 2, 3, 4, 5},
		IncludedInCommitment: nil,
	}
)

type TransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
	tree    *StateTree
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
	s.tree = NewStateTree(s.storage)
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) Test_AddTransaction_AddAndRetrieve() {
	err := s.storage.AddTransaction(&tx)
	s.NoError(err)

	res, err := s.storage.GetTransaction(tx.Hash)
	s.NoError(err)

	s.Equal(tx, *res)
}

func (s *TransactionTestSuite) Test_GetTransaction_NonExistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetTransaction(hash)
	s.NoError(err)
	s.Nil(res)
}

func (s *TransactionTestSuite) Test_GetPendingTransactions_AddAndRetrieve() {
	commitment := &models.Commitment{}
	id, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	tx2 := tx
	tx2.Hash = utils.RandomHash()
	tx3 := tx
	tx3.Hash = utils.RandomHash()
	tx3.IncludedInCommitment = id
	tx4 := tx
	tx4.Hash = utils.RandomHash()
	tx4.ErrorMessage = ref.String("A very boring error message")

	for _, tx := range []*models.Transaction{&tx, &tx2, &tx3, &tx4} {
		err = s.storage.AddTransaction(tx)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingTransactions()
	s.NoError(err)

	s.Equal([]models.Transaction{tx, tx2}, res)
}

func (s *TransactionTestSuite) Test_GetUserTransactions() {
	tx1 := tx
	tx1.Hash = utils.RandomHash()
	tx1.FromIndex = 1
	tx2 := tx
	tx2.Hash = utils.RandomHash()
	tx2.FromIndex = 2
	tx3 := tx
	tx3.Hash = utils.RandomHash()
	tx3.FromIndex = 1

	err := s.storage.AddTransaction(&tx1)
	s.NoError(err)
	err = s.storage.AddTransaction(&tx2)
	s.NoError(err)
	err = s.storage.AddTransaction(&tx3)
	s.NoError(err)

	userTransactions, err := s.storage.GetUserTransactions(models.MakeUint256(1))
	s.NoError(err)

	s.Len(userTransactions, 2)
	s.Contains(userTransactions, tx1)
	s.Contains(userTransactions, tx3)
}

func (s *TransactionTestSuite) Test_SetTransactionError() {
	err := s.storage.AddTransaction(&tx)
	s.NoError(err)

	errorMessage := ref.String("Quack")

	err = s.storage.SetTransactionError(tx.Hash, *errorMessage)
	s.NoError(err)

	res, err := s.storage.GetTransaction(tx.Hash)
	s.NoError(err)

	s.Equal(errorMessage, res.ErrorMessage)
}

func (s *TransactionTestSuite) Test_GetTransactionsByPublicKey() {
	accounts := []models.Account{
		{
			AccountIndex: 1,
			PublicKey:    models.PublicKey{1, 2, 3},
		},
		{
			AccountIndex: 3,
			PublicKey:    models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			AccountIndex: accounts[0].AccountIndex,
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
			AccountIndex: accounts[0].AccountIndex,
			TokenIndex:   models.MakeUint256(25),
			Balance:      models.MakeUint256(1),
			Nonce:        models.MakeUint256(73),
		},
		{
			AccountIndex: accounts[1].AccountIndex,
			TokenIndex:   models.MakeUint256(30),
			Balance:      models.MakeUint256(50),
			Nonce:        models.MakeUint256(71),
		},
	}

	for i := range userStates {
		err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	tx1 := tx
	tx1.Hash = utils.RandomHash()
	tx1.FromIndex = 0
	tx2 := tx
	tx2.Hash = utils.RandomHash()
	tx2.FromIndex = 1
	tx3 := tx
	tx3.Hash = utils.RandomHash()
	tx3.FromIndex = 2
	tx4 := tx
	tx4.Hash = utils.RandomHash()
	tx4.FromIndex = 3

	err := s.storage.AddTransaction(&tx1)
	s.NoError(err)
	err = s.storage.AddTransaction(&tx2)
	s.NoError(err)
	err = s.storage.AddTransaction(&tx3)
	s.NoError(err)
	err = s.storage.AddTransaction(&tx4)
	s.NoError(err)

	userTransactions, err := s.storage.GetTransactionsByPublicKey(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(userTransactions, 3)
	s.Contains(userTransactions, tx1)
	s.Contains(userTransactions, tx3)
	s.Contains(userTransactions, tx4)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
