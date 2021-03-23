package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var cfg = config.RollupConfig{
	FeeReceiverIndex: 3,
	TxsPerCommitment: 32,
}

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
		AccountIndex: 1,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	receiverState := models.UserState{
		AccountIndex: 2,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(0),
		Nonce:        models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		AccountIndex: 3,
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

func (s *ApplyTransactionsTestSuite) Test_ApplyTransactions_AllValid() {
	txs := generateValidTransactions(10)

	validTxs, err := ApplyTransactions(s.storage, txs, &cfg)
	s.NoError(err)

	s.Len(validTxs, 10)
}

func (s *ApplyTransactionsTestSuite) Test_ApplyTransactions_SomeValid() {
	txs := generateValidTransactions(10)
	txs = append(txs, generateInvalidTransactions(10)...)

	validTxs, err := ApplyTransactions(s.storage, txs, &cfg)
	s.NoError(err)

	s.Len(validTxs, 10)
}

func (s *ApplyTransactionsTestSuite) Test_ApplyTransactions_MoreThan32() {
	txs := generateValidTransactions(60)

	validTxs, err := ApplyTransactions(s.storage, txs, &cfg)
	s.NoError(err)

	s.Len(validTxs, 32)

	state, _ := s.tree.Leaf(1)
	s.Equal(models.MakeUint256(32), state.Nonce)
}

func TestApplyTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransactionsTestSuite))
}

func generateValidTransactions(txAmount int) []models.Transaction {
	txs := make([]models.Transaction, 0, txAmount)
	for i := 0; i < txAmount; i++ {
		tx := models.Transaction{
			Hash:      utils.RandomHash(),
			FromIndex: models.MakeUint256(1),
			ToIndex:   models.MakeUint256(2),
			Amount:    models.MakeUint256(1),
			Fee:       models.MakeUint256(1),
			Nonce:     models.MakeUint256(int64(i)),
		}
		txs = append(txs, tx)
	}
	return txs
}

func generateInvalidTransactions(txAmount int) []models.Transaction {
	txs := make([]models.Transaction, 0, txAmount)
	for i := 0; i < txAmount; i++ {
		tx := models.Transaction{
			FromIndex: models.MakeUint256(1),
			ToIndex:   models.MakeUint256(2),
			Amount:    models.MakeUint256(1),
			Fee:       models.MakeUint256(1),
			Nonce:     models.MakeUint256(0),
		}
		txs = append(txs, tx)
	}
	return txs
}
