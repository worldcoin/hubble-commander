package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	bh "github.com/timshannon/badgerhold/v4"
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
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)
	err = s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)
	err = s.storage.AddMassMigration(&massMigration)
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
		err := s.storage.AddTransfer(&txs[i])
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
	err = s.storage.AddTxCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInCommitment := transfer
	transferInCommitment.Hash = common.Hash{5, 5, 5}
	transferInCommitment.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddTransfer(&transferInCommitment)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer)
	s.NoError(err)

	c2t := create2Transfer
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	mm := massMigration
	mm.Hash = common.Hash{6, 7, 8}
	mm.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddMassMigration(&mm)
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

func (s *StoredTransactionTestSuite) TestStoredTxReceipt_CommitmentID_IndexWorks() {
	zeroID := models.CommitmentID{BatchID: models.MakeUint256(0), IndexInBatch: 0}
	id1 := models.CommitmentID{BatchID: models.MakeUint256(1), IndexInBatch: 0}
	id2 := models.CommitmentID{BatchID: models.MakeUint256(2), IndexInBatch: 0}
	s.addStoredTxReceipt(nil, &id1)
	s.addStoredTxReceipt(nil, &id2)
	s.addStoredTxReceipt(nil, &id1)

	indexValues := s.getCommitmentIDIndexValues()
	s.Len(indexValues, 3)
	s.Len(indexValues[zeroID], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[id1], 2)
	s.Len(indexValues[id2], 1)
}

func (s *StoredTransactionTestSuite) TestStoredTxReceipt_CommitmentID_ValuesWithThisFieldSetToNilAreNotIndexed() {
	zeroID := models.CommitmentID{BatchID: models.MakeUint256(0), IndexInBatch: 0}
	s.addStoredTxReceipt(nil, nil)

	indexValues := s.getCommitmentIDIndexValues()
	s.Len(indexValues, 1)
	s.Len(indexValues[zeroID], 0) // value set due to index initialization, see NewTransactionStorage
}

// This test checks an edge case that we introduced by indexing CommitmentID field which can be nil.
// See: NewTransactionStorage
func (s *StoredTransactionTestSuite) TestStoredTxReceipt_CommitmentID_FindUsingIndexWorksWhenThereAreOnlyValuesWithThisFieldSetToNil() {
	err := s.storage.addStoredTxReceipt(&stored.TxReceipt{
		Hash:         utils.RandomHash(),
		CommitmentID: nil, // nil values are not indexed
	})
	s.NoError(err)

	receipts := make([]stored.TxReceipt, 0, 1)
	err = s.storage.database.Badger.Find(
		&receipts,
		bh.Where("CommitmentID").Eq(uint32(1)).Index("CommitmentID"),
	)
	s.NoError(err)
	s.Len(receipts, 0)
}

func (s *StoredTransactionTestSuite) TestStoredTxReceipt_ToStateID_IndexWorks() {
	s.addStoredTxReceipt(ref.Uint32(1), nil)
	s.addStoredTxReceipt(ref.Uint32(2), nil)
	s.addStoredTxReceipt(ref.Uint32(1), nil)

	indexValues := s.getToStateIDIndexValues(stored.TxReceiptName)
	s.Len(indexValues, 3)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[1], 2)
	s.Len(indexValues[2], 1)
}

func (s *StoredTransactionTestSuite) TestStoredTxReceipt_ToStateID_ValuesWithThisFieldSetToNilAreNotIndexed() {
	s.addStoredTxReceipt(nil, nil)

	indexValues := s.getToStateIDIndexValues(stored.TxReceiptName)
	s.Len(indexValues, 1)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
}

// This test checks an edge case that we introduced by indexing ToStateID field which can be nil.
// See: NewTransactionStorage
func (s *StoredTransactionTestSuite) TestStoredTxReceipt_ToStateID_FindUsingIndexWorksWhenThereAreOnlyValuesWithThisFieldSetToNil() {
	err := s.storage.addStoredTxReceipt(&stored.TxReceipt{
		Hash:      utils.RandomHash(),
		ToStateID: nil, // nil values are not indexed
	})
	s.NoError(err)

	receipts := make([]stored.TxReceipt, 0, 1)
	err = s.storage.database.Badger.Find(
		&receipts,
		bh.Where("ToStateID").Eq(uint32(1)).Index("ToStateID"),
	)
	s.NoError(err)
	s.Len(receipts, 0)
}

func (s *StoredTransactionTestSuite) addTransfersInCommitment(batchID *models.Uint256, transfers []models.Transfer) {
	for i := range transfers {
		transfers[i].CommitmentID = &models.CommitmentID{
			BatchID:      *batchID,
			IndexInBatch: 0,
		}
		err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *StoredTransactionTestSuite) addStoredTxReceipt(toStateID *uint32, commitmentID *models.CommitmentID) {
	receipt := &stored.TxReceipt{
		Hash:         utils.RandomHash(),
		ToStateID:    toStateID,
		CommitmentID: commitmentID,
	}
	err := s.storage.addStoredTxReceipt(receipt)
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) getToStateIDIndexValues(typeName []byte) map[uint32]bh.KeyList {
	indexValues := make(map[uint32]bh.KeyList)

	s.iterateIndex(typeName, "ToStateID", func(encodedKey []byte, keyList bh.KeyList) {
		var toStateID uint32
		err := db.Decode(encodedKey, &toStateID)
		s.NoError(err)

		indexValues[toStateID] = keyList
	})

	return indexValues
}

func (s *StoredTransactionTestSuite) getCommitmentIDIndexValues() map[models.CommitmentID]bh.KeyList {
	indexValues := make(map[models.CommitmentID]bh.KeyList)

	s.iterateIndex(stored.TxReceiptName, "CommitmentID", func(encodedKey []byte, keyList bh.KeyList) {
		var commitmentID models.CommitmentID
		err := db.Decode(encodedKey, &commitmentID)
		s.NoError(err)

		indexValues[commitmentID] = keyList
	})

	return indexValues
}

func (s *StoredTransactionTestSuite) iterateIndex(
	typeName []byte,
	indexName string,
	handleIndex func(encodedKey []byte, keyList bh.KeyList),
) {
	testutils.IterateIndex(s.Assertions, s.storage.database.Badger, typeName, indexName, handleIndex)
}

func TestStoredTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(StoredTransactionTestSuite))
}
