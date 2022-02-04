package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoredTransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StoredTransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StoredTransactionTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) TestSetTransactionErrors() {
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

func (s *StoredTransactionTestSuite) TestMarkTransactionsAsPending() {
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

func (s *StoredTransactionTestSuite) TestGetTransactionCount() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
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

	count, err := s.storage.GetTransactionCount()
	s.NoError(err)
	s.Equal(3, *count)
}

func (s *StoredTransactionTestSuite) TestGetTransactionCount_NoTransactions() {
	count, err := s.storage.GetTransactionCount()
	s.NoError(err)
	s.Equal(0, *count)
}

func (s *StoredTransactionTestSuite) TestGetTransactionHashesByBatchIDs() {
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

func (s *StoredTransactionTestSuite) TestGetTransactionHashesByBatchIDs_NoTransactions() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = transfer
	transfers[1] = transfer
	transfers[1].Hash = utils.RandomHash()
	s.addTransfersInCommitment(models.NewUint256(1), transfers)

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(models.MakeUint256(2))
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(hashes)
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_Transfers() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)

	txs, err := s.storage.GetPendingTransactions(txtype.Transfer)
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(transfer, *txs.At(0).ToTransfer())
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_NoTransfers() {
	txs, err := s.storage.GetPendingTransactions(txtype.Transfer)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_Create2Transfers() {
	err := s.storage.AddTransaction(&create2Transfer)
	s.NoError(err)

	txs, err := s.storage.GetPendingTransactions(txtype.Create2Transfer)
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(create2Transfer, *txs.At(0).ToCreate2Transfer())
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_NoCreate2Transfers() {
	txs, err := s.storage.GetPendingTransactions(txtype.Create2Transfer)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_MassMigrations() {
	err := s.storage.AddTransaction(&massMigration)
	s.NoError(err)

	txs, err := s.storage.GetPendingTransactions(txtype.MassMigration)
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(massMigration, *txs.At(0).ToMassMigration())
}

func (s *StoredTransactionTestSuite) TestGetPendingTransactions_NoMassMigrations() {
	txs, err := s.storage.GetPendingTransactions(txtype.MassMigration)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *StoredTransactionTestSuite) TestGetAllFailedTransactions() {
	failedTransfer := transfer
	failedTransfer.ErrorMessage = ref.String("quacked transfer")
	err := s.storage.AddTransaction(&failedTransfer)
	s.NoError(err)

	failedC2T := create2Transfer
	failedC2T.ErrorMessage = ref.String("quacked create2transfer")
	err = s.storage.AddTransaction(&failedC2T)
	s.NoError(err)

	failedMM := massMigration
	failedMM.ErrorMessage = ref.String("quacked migration")
	err = s.storage.AddTransaction(&failedMM)
	s.NoError(err)

	txs, err := s.storage.GetAllFailedTxs()
	s.NoError(err)
	s.Len(txs, 3)
}

func (s *StoredTransactionTestSuite) TestGetAllFailedTransactions_NoTransactions() {
	txs, err := s.storage.GetAllFailedTxs()
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *StoredTransactionTestSuite) addTransfersInCommitment(batchID *models.Uint256, transfers []models.Transfer) {
	for i := range transfers {
		transfers[i].CommitmentID = &models.CommitmentID{
			BatchID:      *batchID,
			IndexInBatch: 0,
		}
		err := s.storage.AddTransaction(&transfers[i])
		s.NoError(err)
	}
}

func TestStoredTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(StoredTransactionTestSuite))
}
