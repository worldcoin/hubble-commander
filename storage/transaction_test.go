package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	batch   *models.Batch
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)

	s.batch = &models.Batch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		MinedTime:       &models.Timestamp{Time: time.Unix(140, 0).UTC()},
	}

	err = s.storage.AddBatch(s.batch)
	s.NoError(err)
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_WithoutBatch() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: transfer.Copy(),
	}

	res, err := s.storage.GetTransactionWithBatchDetails(transfer.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_Transfer() {
	transferInBatch := transfer
	transferInBatch.CommitmentID = &models.CommitmentID{
		BatchID: s.batch.ID,
	}
	err := s.storage.AddTransaction(&transferInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: transferInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		MinedTime:   s.batch.MinedTime,
	}
	res, err := s.storage.GetTransactionWithBatchDetails(transferInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_Create2Transfer() {
	create2TransferInBatch := create2Transfer
	create2TransferInBatch.CommitmentID = &models.CommitmentID{
		BatchID: s.batch.ID,
	}
	err := s.storage.AddTransaction(&create2TransferInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: create2TransferInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		MinedTime:   s.batch.MinedTime,
	}
	res, err := s.storage.GetTransactionWithBatchDetails(create2TransferInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_MassMigration() {
	massMigrationInBatch := massMigration
	massMigrationInBatch.CommitmentID = &models.CommitmentID{
		BatchID: s.batch.ID,
	}
	err := s.storage.AddTransaction(&massMigrationInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: massMigrationInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		MinedTime:   s.batch.MinedTime,
	}
	res, err := s.storage.GetTransactionWithBatchDetails(massMigrationInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransactionTestSuite) TestReplaceFailedTransaction_UpdatesTx() {
	failedTx := transfer
	failedTx.ErrorMessage = ref.String("some message")
	err := s.storage.AddTransaction(&failedTx)
	s.NoError(err)

	updatedTx := transfer
	updatedTx.SetReceiveTime()
	err = s.storage.ReplaceFailedTransaction(&updatedTx)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Equal(updatedTx, *res)
}

func (s *TransactionTestSuite) TestReplaceFailedTransaction_DoesNotUpdateBatchedTx() {
	tx := create2Transfer
	tx.CommitmentID = &models.CommitmentID{
		BatchID: models.MakeUint256(1),
	}
	err := s.storage.AddTransaction(&tx)
	s.NoError(err)

	updatedTx := create2Transfer
	updatedTx.SetReceiveTime()
	err = s.storage.ReplaceFailedTransaction(&updatedTx)
	s.ErrorIs(err, ErrAlreadyMinedTransaction)
}

func (s *TransactionTestSuite) TestReplaceFailedTransaction_DoesNotUpdatePendingTx() {
	err := s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	updatedTx := massMigration
	updatedTx.SetReceiveTime()
	err = s.storage.ReplaceFailedTransaction(&updatedTx)
	s.ErrorIs(err, NewNotFoundError("FailedTx"))
}

func (s *TransactionTestSuite) TestReplacePendingTransaction() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)

	updatedTx := transfer
	updatedTx.Hash = utils.RandomHash()
	err = s.storage.ReplacePendingTransaction(&transfer.Hash, &updatedTx)
	s.NoError(err)

	_, err = s.storage.GetTransfer(transfer.Hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))

	tx, err := s.storage.GetTransfer(updatedTx.Hash)
	s.NoError(err)
	s.Equal(updatedTx, *tx)
}

func (s *TransactionTestSuite) TestReplacePendingTransaction_NoPendingTransaction() {
	updatedTx := transfer
	updatedTx.Hash = utils.RandomHash()
	err := s.storage.ReplacePendingTransaction(&transfer.Hash, &updatedTx)
	s.ErrorIs(err, NewNotFoundError("transaction"))
}

func (s *TransactionTestSuite) TestGetTransactionsByCommitmentID() {
	transfer1 := transfer
	transfer1.CommitmentID = &txCommitment.ID
	err := s.storage.AddTransaction(&transfer1)
	s.NoError(err)

	otherCommitmentID := txCommitment.ID
	otherCommitmentID.IndexInBatch += 1
	transfer2 := create2Transfer
	transfer2.Hash = utils.RandomHash()
	transfer2.CommitmentID = &otherCommitmentID
	err = s.storage.AddTransaction(&transfer2)
	s.NoError(err)

	transfers, err := s.storage.GetTransactionsByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 1)
	s.Equal(transfers.At(0), &transfer1)
}

func (s *TransactionTestSuite) TestGetTransactionsByCommitmentID_NoTransactions() {
	transfers, err := s.storage.GetTransactionsByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *TransactionTestSuite) TestBatchUpsertTransaction() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)

	txBeforeUpsert, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Nil(txBeforeUpsert.CommitmentID)

	txBeforeUpsert.CommitmentID = &models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}

	err = s.storage.BatchUpsertTransaction(models.GenericArray{txBeforeUpsert})
	s.NoError(err)

	txAfterUpsert, err := s.storage.GetTransfer(txBeforeUpsert.Hash)
	s.NoError(err)
	s.Equal(txBeforeUpsert, txAfterUpsert)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
