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

func (s *TransactionTestSuite) TestUpdateTransaction_UpdatesTx() {
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

func (s *TransactionTestSuite) TestUpdateTransaction_DoesNotUpdateMinedTx() {
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

func (s *TransactionTestSuite) TestUpdateTransaction_DoesNotUpdatePendingTx() {
	err := s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	updatedTx := massMigration
	updatedTx.SetReceiveTime()
	err = s.storage.ReplaceFailedTransaction(&updatedTx)
	s.ErrorIs(err, NewNotFoundError("FailedTx"))
}

func (s *TransactionTestSuite) populatePendingTransactions() models.GenericTransactionArray {
	transfers := make([]models.Transfer, 4)
	for i := range transfers {
		transfers[i] = transfer
		transfers[i].Hash = utils.RandomHash()
	}
	transfers[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(3)}
	transfers[3].ErrorMessage = ref.String("A very boring error message")

	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)

	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			Type: batchtype.Transfer,
		},
	}
	err = s.storage.AddCommitment(commitment)
	s.NoError(err)

	create2Transfer2 := create2Transfer
	create2Transfer2.Hash = utils.RandomHash()
	create2Transfer3 := create2Transfer
	create2Transfer3.Hash = utils.RandomHash()
	create2Transfer3.CommitmentID = &commitment.ID
	create2Transfer4 := create2Transfer
	create2Transfer4.Hash = utils.RandomHash()
	create2Transfer4.ErrorMessage = ref.String("A very boring error message")

	create2transfers := []models.Create2Transfer{
		create2Transfer,
		create2Transfer2,
		create2Transfer3,
		create2Transfer4,
	}

	err = s.storage.BatchAddCreate2Transfer(create2transfers)
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
	result = result.Append(models.MakeCreate2TransferArray(create2transfers...))
	result = result.Append(models.MakeMassMigrationArray(massMigrations...))
	result = result.Append(models.MakeTransferArray(transfers...))
	return result
}

func (s *TransactionTestSuite) TestGetPendingTransactions_GetCreate2Transfers() {
	transactions := s.populatePendingTransactions()
	create2transfers := transactions.ToCreate2TransferArray()

	res, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, create2transfers[0])
	s.Contains(res, create2transfers[1])
}

func (s *TransactionTestSuite) TestGetPendingTransactions_GetMassMigrations() {
	transactions := s.populatePendingTransactions()
	massMigrations := transactions.ToMassMigrationArray()

	res, err := s.storage.GetPendingMassMigrations()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, massMigrations[0])
	s.Contains(res, massMigrations[1])
}

func (s *TransactionTestSuite) TestGetPendingTransactions_GetTransfers() {
	transactions := s.populatePendingTransactions()
	transfers := transactions.ToTransferArray()

	res, err := s.storage.GetPendingTransfers()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, transfers[0])
	s.Contains(res, transfers[1])
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
