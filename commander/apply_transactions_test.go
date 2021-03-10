package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyTransactionsTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	tree    *storage.StateTree
}

func (s *ApplyTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransactionsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.tree = storage.NewStateTree(s.storage)

	senderState := models.UserState{
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	receiverState := models.UserState{
		AccountIndex: models.MakeUint256(2),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(0),
		Nonce:        models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		AccountIndex: models.MakeUint256(3),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(1000),
		Nonce:        models.MakeUint256(0),
	}

	err = s.tree.Set(1, &senderState)
	s.NoError(err)
	err = s.tree.Set(2, &receiverState)
	s.NoError(err)
	err = s.tree.Set(3, &feeReceiverState)
	s.NoError(err)
}

func (s *ApplyTransactionsTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *ApplyTransactionsTestSuite) Test_ApplyTransfers_AllValid() {
	transactions := generateValidTransactions(10)

	validTransactions, err := ApplyTransactions(s.tree, transactions, uint32(3))
	s.NoError(err)

	s.Len(validTransactions, 10)
}

func (s *ApplyTransactionsTestSuite) Test_ApplyTransfers_SomeValid() {
	transactions := generateValidTransactions(10)
	transactions = append(transactions, generateInvalidTransactions(10)...)

	validTransactions, err := ApplyTransactions(s.tree, transactions, uint32(3))
	s.NoError(err)

	s.Len(validTransactions, 10)
}

func (s *ApplyTransactionsTestSuite) Test_ApplyTransfers_MoreThan32() {
	transactions := generateValidTransactions(60)

	validTransactions, err := ApplyTransactions(s.tree, transactions, uint32(3))
	s.NoError(err)

	s.Len(validTransactions, 32)

	state, _ := s.tree.Leaf(1)
	s.Equal(models.MakeUint256(32), state.Nonce)
}

func TestApplyTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransactionsTestSuite))
}

func generateValidTransactions(txAmount int) []models.Transaction {
	transactions := make([]models.Transaction, 0, txAmount)
	for i := 0; i < txAmount; i++ {
		transaction := models.Transaction{
			FromIndex: models.MakeUint256(1),
			ToIndex:   models.MakeUint256(2),
			Amount:    models.MakeUint256(1),
			Fee:       models.MakeUint256(1),
			Nonce:     models.MakeUint256(int64(i)),
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

func generateInvalidTransactions(txAmount int) []models.Transaction {
	transactions := make([]models.Transaction, 0, txAmount)
	for i := 0; i < txAmount; i++ {
		transaction := models.Transaction{
			FromIndex: models.MakeUint256(1),
			ToIndex:   models.MakeUint256(2),
			Amount:    models.MakeUint256(1),
			Fee:       models.MakeUint256(1),
			Nonce:     models.MakeUint256(0),
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}
