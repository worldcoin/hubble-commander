package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	bh "github.com/timshannon/badgerhold/v4"
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

func (s *TransactionTestSuite) TestAddPendingTransactions_AddTxs() {
	txs := models.MakeGenericArray(
		transfer.ToTransfer(),
		massMigration.ToMassMigration(),
		create2Transfer.ToCreate2Transfer(),
	)
	err := s.storage.AddPendingTransactions(txs)
	s.NoError(err)

	gotTxs, err := s.storage.GetAllPendingTransactions()
	s.NoError(err)
	s.Equal(txs.Len(), gotTxs.Len())
}

func (s *TransactionTestSuite) TestAddFailedTransactions_NoTransactions() {
	arr := models.MakeTransferArray()
	err := s.storage.AddFailedTransactions(arr)
	s.NoError(err)
}

func (s *TransactionTestSuite) TestAddFailedTransactions_TransactionWithError() {
	arr := models.MakeTransferArray(transfer)
	err := s.storage.AddFailedTransactions(arr)
	s.NoError(err)

	err = s.storage.AddFailedTransactions(arr)
	s.ErrorIs(err, bh.ErrKeyExists)
}

func (s *TransactionTestSuite) TestSetTransactionErrors() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)
	err = s.storage.AddTransaction(&create2Transfer)
	s.NoError(err)
	err = s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	transferError := models.TxError{
		TxHash:       transfer.Hash,
		ErrorMessage: "Quack",
	}

	c2tError := models.TxError{
		TxHash:       create2Transfer.Hash,
		ErrorMessage: "C2T Quack",
	}

	mmError := models.TxError{
		TxHash:       massMigration.Hash,
		ErrorMessage: "MM Quack",
	}

	err = s.storage.SetTransactionErrors(transferError, c2tError, mmError)
	s.NoError(err)

	storedTransfer, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Equal(transferError.ErrorMessage, *storedTransfer.ErrorMessage)

	storedC2T, err := s.storage.GetCreate2Transfer(create2Transfer.Hash)
	s.NoError(err)
	s.Equal(c2tError.ErrorMessage, *storedC2T.ErrorMessage)

	storedMM, err := s.storage.GetMassMigration(massMigration.Hash)
	s.NoError(err)
	s.Equal(mmError.ErrorMessage, *storedMM.ErrorMessage)
}

func (s *TransactionTestSuite) TestMarkTransactionsAsPending_BatchedTxs() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transfer
		txs[i].Hash = utils.RandomHash()
		txs[i].CommitmentID = &models.CommitmentID{
			BatchID:      models.MakeUint256(5),
			IndexInBatch: 3,
		}
		err := s.storage.AddTransaction(&txs[i])
		s.NoError(err)
	}

	err := s.storage.MarkTransactionsAsPending([]common.Hash{txs[0].Hash, txs[1].Hash})
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		s.Nil(tx.CommitmentID)
	}
}

func (s *TransactionTestSuite) TestMarkTransactionsAsPending_FailedTxs() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transfer
		txs[i].ErrorMessage = ref.String("error message")
		txs[i].Hash = utils.RandomHash()
		txs[i].CommitmentID = &models.CommitmentID{
			BatchID:      models.MakeUint256(5),
			IndexInBatch: 3,
		}
		err := s.storage.AddTransaction(&txs[i])
		s.NoError(err)
	}

	err := s.storage.MarkTransactionsAsPending([]common.Hash{txs[0].Hash, txs[1].Hash})
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		s.Nil(tx.CommitmentID)
	}
}

func (s *TransactionTestSuite) TestMarkTransactionsAsPending_NotFound() {
	err := s.storage.MarkTransactionsAsPending([]common.Hash{utils.RandomHash()})
	s.Error(err)
}

func (s *TransactionTestSuite) TestGetTransactionCount() {
	batch := &models.Batch{
		ID:                models.MakeUint256(0),
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := txCommitment
	commitmentInBatch.ID.BatchID = batch.ID
	err = s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInCommitment := transfer
	transferInCommitment.Hash = common.Hash{5, 5, 5}
	transferInCommitment.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddTransaction(&transferInCommitment)
	s.NoError(err)
	err = s.storage.AddTransaction(&transfer)
	s.NoError(err)

	c2t := create2Transfer
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddTransaction(&c2t)
	s.NoError(err)

	mm := massMigration
	mm.Hash = common.Hash{6, 7, 8}
	mm.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddTransaction(&mm)
	s.NoError(err)

	storageCount, err := s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(3, *storageCount)
	count := s.storage.GetTransactionCount()
	s.EqualValues(3, count)

	err = s.storage.MarkTransactionsAsPending([]common.Hash{transferInCommitment.Hash})
	s.NoError(err)
	storageCount, err = s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(2, *storageCount)
	count = s.storage.GetTransactionCount()
	s.EqualValues(2, count)

	err = s.storage.MarkTransfersAsIncluded([]models.Transfer{transferInCommitment}, &commitmentInBatch.ID)
	s.NoError(err)
	storageCount, err = s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(3, *storageCount)
	count = s.storage.GetTransactionCount()
	s.EqualValues(3, count)
}

func (s *TransactionTestSuite) TestGetTransactionCount_IncrementsTxCountOnStorageCopy() {
	transactionStorageCopy := s.storage.TransactionStorage.copyWithNewDatabase(s.storage.database)
	transactionStorageCopy.incrementTransactionCount()

	count := s.storage.GetTransactionCount()
	s.EqualValues(1, count)
}

func (s *TransactionTestSuite) TestGetTransactionCount_NoTransactions() {
	count := s.storage.GetTransactionCount()
	s.EqualValues(0, count)
}

func (s *TransactionTestSuite) TestGetTransactionHashesByBatchIDs() {
	batchIDs := []models.Uint256{models.MakeUint256(1), models.MakeUint256(2)}
	expectedHashes := make([]common.Hash, 0, 4)
	for i := range batchIDs {
		transfers := make([]models.Transfer, 2)
		transfers[0] = transfer
		transfers[0].Hash = utils.RandomHash()
		transfers[1] = transfer
		transfers[1].Hash = utils.RandomHash()
		s.addTransfersInCommitment(&batchIDs[i], transfers)
		expectedHashes = append(expectedHashes, transfers[0].Hash, transfers[1].Hash)
	}

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(batchIDs...)
	s.NoError(err)
	s.Len(hashes, 4)
	for i := range expectedHashes {
		s.Contains(hashes, expectedHashes[i])
	}
}

func (s *TransactionTestSuite) TestGetTransactionHashesByBatchIDs_NoTransactions() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = transfer
	transfers[1] = transfer
	transfers[1].Hash = utils.RandomHash()
	s.addTransfersInCommitment(models.NewUint256(1), transfers)

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(models.MakeUint256(2))
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(hashes)
}

func (s *TransactionTestSuite) TestGetPendingTransactions_Transfers() {
	transactions := s.populateTransactions()
	transfers := transactions.ToTransferArray()

	res, err := s.storage.GetPendingTransactions(txtype.Transfer)
	s.NoError(err)
	s.Len(res, 2)

	transferArray := res.ToTransferArray()
	s.Contains(transferArray, transfers[0])
	s.Contains(transferArray, transfers[1])
}

func (s *TransactionTestSuite) TestGetPendingTransactions_NoTransfers() {
	txs, err := s.storage.GetPendingTransactions(txtype.Transfer)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *TransactionTestSuite) TestGetPendingTransactions_Create2Transfers() {
	transactions := s.populateTransactions()
	create2transfers := transactions.ToCreate2TransferArray()

	res, err := s.storage.GetPendingTransactions(txtype.Create2Transfer)
	s.NoError(err)
	s.Len(res, 2)

	c2tArray := res.ToCreate2TransferArray()
	s.Contains(c2tArray, create2transfers[0])
	s.Contains(c2tArray, create2transfers[1])
}

func (s *TransactionTestSuite) TestGetPendingTransactions_NoCreate2Transfers() {
	txs, err := s.storage.GetPendingTransactions(txtype.Create2Transfer)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *TransactionTestSuite) TestGetPendingTransactions_MassMigrations() {
	transactions := s.populateTransactions()
	massMigrations := transactions.ToMassMigrationArray()

	res, err := s.storage.GetPendingTransactions(txtype.MassMigration)
	s.NoError(err)
	s.Len(res, 2)

	mmArray := res.ToMassMigrationArray()
	s.Contains(mmArray, massMigrations[0])
	s.Contains(mmArray, massMigrations[1])
}

func (s *TransactionTestSuite) TestGetPendingTransactions_NoMassMigrations() {
	txs, err := s.storage.GetPendingTransactions(txtype.MassMigration)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *TransactionTestSuite) TestGetAllPendingTransactions() {
	txs := s.populateTransactions()
	expectedTxs := make([]models.GenericTransaction, 0, 6)
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)
		base := tx.GetBase()
		if base.CommitmentID == nil && base.ErrorMessage == nil {
			expectedTxs = append(expectedTxs, tx)
		}
	}

	res, err := s.storage.GetAllPendingTransactions()
	s.NoError(err)
	s.Len(res, 6)
	for _, expectedTx := range expectedTxs {
		s.Contains(res, expectedTx)
	}
}

func (s *TransactionTestSuite) TestGetAllPendingTransactions_NoTransactions() {
	txs, err := s.storage.GetAllPendingTransactions()
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *TransactionTestSuite) TestGetAllFailedTransactions() {
	txs := s.populateTransactions()
	expectedTxs := make([]models.GenericTransaction, 0, 3)
	for i := 0; i < txs.Len(); i++ {
		tx := txs.At(i)
		if tx.GetBase().ErrorMessage != nil {
			expectedTxs = append(expectedTxs, tx)
		}
	}

	res, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(res, 3)
	for _, expectedTx := range expectedTxs {
		s.Contains(res, expectedTx)
	}
}

func (s *TransactionTestSuite) TestGetAllFailedTransactions_NoTransactions() {
	txs, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *TransactionTestSuite) populateTransactions() models.GenericTransactionArray {
	transfers := make([]models.Transfer, 4)
	for i := range transfers {
		transfers[i] = transfer
		transfers[i].Hash = utils.RandomHash()
	}
	transfers[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(1)}
	transfers[3].ErrorMessage = ref.String("A very boring error message")

	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)

	create2Transfers := make([]models.Create2Transfer, 4)
	for i := range create2Transfers {
		create2Transfers[i] = create2Transfer
		create2Transfers[i].Hash = utils.RandomHash()
	}

	create2Transfers[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(2)}
	create2Transfers[3].ErrorMessage = ref.String("A very boring error message")

	err = s.storage.BatchAddCreate2Transfer(create2Transfers)
	s.NoError(err)

	massMigrations := make([]models.MassMigration, 4)
	for i := range massMigrations {
		massMigrations[i] = massMigration
		massMigrations[i].Hash = utils.RandomHash()
	}
	massMigrations[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(3)}
	massMigrations[3].ErrorMessage = ref.String("A very boring error message")

	err = s.storage.BatchAddMassMigration(massMigrations)
	s.NoError(err)

	var result models.GenericTransactionArray
	result = models.MakeGenericArray()
	result = result.Append(models.MakeTransferArray(transfers...))
	result = result.Append(models.MakeCreate2TransferArray(create2Transfers...))
	result = result.Append(models.MakeMassMigrationArray(massMigrations...))
	return result
}

func (s *TransactionTestSuite) addTransfersInCommitment(batchID *models.Uint256, transfers []models.Transfer) {
	for i := range transfers {
		transfers[i].CommitmentID = &models.CommitmentID{
			BatchID:      *batchID,
			IndexInBatch: 0,
		}
		err := s.storage.AddTransaction(&transfers[i])
		s.NoError(err)
	}
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
