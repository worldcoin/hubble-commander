package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	transferTransaction = models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			FromStateID:          1,
			Amount:               models.MakeUint256(1000),
			Fee:                  models.MakeUint256(100),
			Nonce:                models.MakeUint256(0),
			Signature:            models.MakeRandomSignature(),
			IncludedInCommitment: nil,
		},
		ToStateID: 2,
	}
)

type TransactionBaseTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *postgres.TestDB
}

func (s *TransactionBaseTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionBaseTestSuite) SetupTest() {
	testDB, err := postgres.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *TransactionBaseTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransactionBaseTestSuite) TestSetTransactionError() {
	err := s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)

	errorMessage := ref.String("Quack")

	err = s.storage.SetTransactionError(transferTransaction.Hash, *errorMessage)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transferTransaction.Hash)
	s.NoError(err)

	s.Equal(errorMessage, res.ErrorMessage)
}

func (s *TransactionBaseTestSuite) TestGetLatestTransactionNonce() {
	account := models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(3)
	tx2 := transferTransaction
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(5)
	tx3 := transferTransaction
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(1)

	err = s.storage.AddTransfer(&tx1)
	s.NoError(err)
	err = s.storage.AddTransfer(&tx2)
	s.NoError(err)
	err = s.storage.AddTransfer(&tx3)
	s.NoError(err)

	userTransactions, err := s.storage.GetLatestTransactionNonce(account.PubKeyID)
	s.NoError(err)
	s.Equal(models.NewUint256(5), userTransactions)
}

func TestTransactionBaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionBaseTestSuite))
}
