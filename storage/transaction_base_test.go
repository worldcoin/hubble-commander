package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	transferTransaction = models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: 2,
	}
	create2TransferTransaction = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToPublicKey: models.PublicKey{1, 2, 3},
	}
)

type TransactionBaseTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *TransactionBaseTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionBaseTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
}

func (s *TransactionBaseTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransactionBaseTestSuite) TestSetTransactionError() {
	_, err := s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)

	errorMessage := ref.String("Quack")

	err = s.storage.SetTransactionError(transferTransaction.Hash, *errorMessage)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transferTransaction.Hash)
	s.NoError(err)

	s.Equal(errorMessage, res.ErrorMessage)
}

func (s *TransactionBaseTestSuite) TestGetLatestTransactionNonce() {
	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(3)

	tx2 := create2TransferTransaction
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(5)

	tx3 := transferTransaction
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(1)

	_, err := s.storage.AddTransfer(&tx1)
	s.NoError(err)
	_, err = s.storage.AddCreate2Transfer(&tx2)
	s.NoError(err)
	_, err = s.storage.AddTransfer(&tx3)
	s.NoError(err)

	userTransactions, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(5), userTransactions)
}

func (s *TransactionBaseTestSuite) TestGetLatestTransactionNonce_DisregardsFailedTransactions() {
	s.T().SkipNow() // TODO fix
	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(1)

	tx2 := transferTransaction
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(2)

	tx3 := transferTransaction
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(3)
	tx3.ErrorMessage = ref.String("error")

	txs := []models.Transfer{tx1, tx2, tx3}
	for i := range txs {
		_, err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	userTransactions, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(2), userTransactions)
}

func (s *TransactionBaseTestSuite) TestBatchMarkTransactionAsIncluded() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transferTransaction
		txs[i].Hash = utils.RandomHash()
		_, err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.BatchMarkTransactionAsIncluded([]common.Hash{txs[0].Hash, txs[1].Hash}, commitmentID)
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		s.Equal(commitmentID, tx.IncludedInCommitment)
	}
}

func (s *TransactionBaseTestSuite) TestGetTransactionCount() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := commitment
	commitmentInBatch.IncludedInBatch = &batch.ID
	commitmentID, err := s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInCommitment := transferTransaction
	transferInCommitment.Hash = common.Hash{5, 5, 5}
	transferInCommitment.IncludedInCommitment = commitmentID
	_, err = s.storage.AddTransfer(&transferInCommitment)
	s.NoError(err)
	_, err = s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)

	c2t := create2Transfer
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.IncludedInCommitment = commitmentID
	_, err = s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	count, err := s.storage.GetTransactionCount()
	s.NoError(err)
	s.Equal(2, *count)
}

func (s *TransactionBaseTestSuite) TestGetTransactionCount_NoTransactions() {
	count, err := s.storage.GetTransactionCount()
	s.NoError(err)
	s.Equal(0, *count)
}

func (s *TransactionBaseTestSuite) TestGetTransactionHashesByBatchIDs() {
	batchIDs := []models.Uint256{models.MakeUint256(1), models.MakeUint256(2)}
	for i := range batchIDs {
		transfers := make([]models.Transfer, 2)
		transfers[0] = transfer
		transfers[0].Hash = utils.RandomHash()
		transfers[1] = transfer
		transfers[1].Hash = utils.RandomHash()
		s.addTransfersInCommitment(&batchIDs[i], transfers)
	}

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	s.NoError(err)
	s.Len(hashes, 4)
}

func (s *TransactionBaseTestSuite) TestGetTransactionHashesByBatchIDs_NoTransactions() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = transfer
	transfers[1] = transfer
	transfers[1].Hash = utils.RandomHash()
	s.addTransfersInCommitment(models.NewUint256(1), transfers)

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(models.MakeUint256(2))
	s.Equal(NewNotFoundError("transaction"), err)
	s.Nil(hashes)
}

func (s *TransactionBaseTestSuite) addTransfersInCommitment(batchID *models.Uint256, transfers []models.Transfer) {
	batch := &models.Batch{
		ID:                *batchID,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := commitment
	commitmentInBatch.IncludedInBatch = &batch.ID
	commitmentID, err := s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	for i := range transfers {
		transfers[i].IncludedInCommitment = commitmentID
		_, err = s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}
}

func TestTransactionBaseTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionBaseTestSuite))
}
