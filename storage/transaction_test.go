package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	tx = models.Transaction{
		Hash:                 common.BigToHash(big.NewInt(1234)),
		FromIndex:            models.MakeUint256(1),
		ToIndex:              models.MakeUint256(2),
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
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
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
	commitmentHash := common.BigToHash(big.NewInt(1234))

	tx2 := tx
	tx2.Hash = common.BigToHash(big.NewInt(2345))
	tx3 := tx
	tx3.Hash = common.BigToHash(big.NewInt(3456))
	tx3.IncludedInCommitment = &commitmentHash
	tx4 := tx
	tx4.Hash = common.BigToHash(big.NewInt(4567))
	tx4.ErrorMessage = ref.String("A very boring error message")

	for _, tx := range []*models.Transaction{&tx, &tx2, &tx3, &tx4} {
		err := s.storage.AddTransaction(tx)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingTransactions()
	s.NoError(err)

	s.Equal([]models.Transaction{tx, tx2}, res)
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

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
