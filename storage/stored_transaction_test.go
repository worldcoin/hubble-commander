package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	bh "github.com/timshannon/badgerhold/v4"
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
			Hash:        common.BigToHash(big.NewInt(5678)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToPublicKey: models.PublicKey{1, 2, 3},
	}
	massMigrationTransaction = models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		SpokeID: models.MakeUint256(5),
	}
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
	err := s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)
	err = s.storage.AddCreate2Transfer(&create2TransferTransaction)
	s.NoError(err)
	err = s.storage.AddMassMigration(&massMigrationTransaction)
	s.NoError(err)

	transferError := models.TxError{
		TxHash:       transferTransaction.Hash,
		ErrorMessage: "Quack",
	}

	c2tError := models.TxError{
		TxHash:       create2TransferTransaction.Hash,
		ErrorMessage: "C2T Quack",
	}

	mmError := models.TxError{
		TxHash:       massMigrationTransaction.Hash,
		ErrorMessage: "MM Quack",
	}

	err = s.storage.SetTransactionErrors(transferError, c2tError, mmError)
	s.NoError(err)

	storedTransfer, err := s.storage.GetTransfer(transferTransaction.Hash)
	s.NoError(err)
	s.Equal(transferError.ErrorMessage, *storedTransfer.ErrorMessage)

	storedC2T, err := s.storage.GetCreate2Transfer(create2TransferTransaction.Hash)
	s.NoError(err)
	s.Equal(c2tError.ErrorMessage, *storedC2T.ErrorMessage)

	storedMM, err := s.storage.GetMassMigration(massMigrationTransaction.Hash)
	s.NoError(err)
	s.Equal(mmError.ErrorMessage, *storedMM.ErrorMessage)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_ReturnsHighestNonceRegardlessOfInsertionOrder() {
	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(3)

	tx2 := transferTransaction
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(5)

	tx3 := transferTransaction
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(1)

	txs := []models.Transfer{tx1, tx2, tx3}
	for i := range txs {
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(5), latestNonce)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_ReturnsHighestNonceRegardlessOfTxType() {
	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(3)

	tx2 := create2TransferTransaction
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(5)

	tx3 := massMigrationTransaction
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(4)

	err := s.storage.AddTransfer(&tx1)
	s.NoError(err)
	err = s.storage.AddCreate2Transfer(&tx2)
	s.NoError(err)
	err = s.storage.AddMassMigration(&tx3)
	s.NoError(err)

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(5), latestNonce)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_DisregardsTransactionsFromOtherStateIDs() {
	tx1 := transferTransaction
	tx1.Hash = utils.RandomHash()
	tx1.Nonce = models.MakeUint256(3)

	tx2 := transferTransaction
	tx2.FromStateID = 10
	tx2.Hash = utils.RandomHash()
	tx2.Nonce = models.MakeUint256(5)

	tx3 := transferTransaction
	tx3.FromStateID = 20
	tx3.Hash = utils.RandomHash()
	tx3.Nonce = models.MakeUint256(7)

	txs := []models.Transfer{tx1, tx2, tx3}
	for i := range txs {
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(3), latestNonce)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_DisregardsFailedTransactions() {
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
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.NoError(err)
	s.Equal(models.NewUint256(2), latestNonce)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_NoTransactionsForGivenStateID() {
	tx1 := transferTransaction
	tx1.FromStateID = 10

	err := s.storage.AddTransfer(&tx1)
	s.NoError(err)

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(latestNonce)
}

func (s *StoredTransactionTestSuite) TestGetLatestTransactionNonce_NoValidTransactionsForGivenStateID() {
	tx1 := transferTransaction
	tx1.ErrorMessage = ref.String("error")

	err := s.storage.AddTransfer(&tx1)
	s.NoError(err)

	latestNonce, err := s.storage.GetLatestTransactionNonce(1)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(latestNonce)
}

func (s *StoredTransactionTestSuite) TestMarkTransactionsAsPending() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transferTransaction
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

	transferInCommitment := transferTransaction
	transferInCommitment.Hash = common.Hash{5, 5, 5}
	transferInCommitment.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddTransfer(&transferInCommitment)
	s.NoError(err)
	err = s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)

	c2t := create2TransferTransaction
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	mm := massMigrationTransaction
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
		transfers[0] = transferTransaction
		transfers[0].Hash = utils.RandomHash()
		transfers[1] = transferTransaction
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
	transfers[0] = transferTransaction
	transfers[1] = transferTransaction
	transfers[1].Hash = utils.RandomHash()
	s.addTransfersInCommitment(models.NewUint256(1), transfers)

	hashes, err := s.storage.GetTransactionHashesByBatchIDs(models.MakeUint256(2))
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(hashes)
}

func (s *StoredTransactionTestSuite) TestStoredTx_ToStateID_IndexWorks() {
	s.addStoredTx(txtype.Transfer, ref.Uint32(1))
	s.addStoredTx(txtype.Transfer, ref.Uint32(2))
	s.addStoredTx(txtype.Transfer, ref.Uint32(1))

	indexValues := s.getToStateIDIndexValues(models.StoredTxName)
	s.Len(indexValues, 3)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[1], 2)
	s.Len(indexValues[2], 1)
}

func (s *StoredTransactionTestSuite) TestStoredTx_ToStateID_ValuesWithoutThisFieldAreNotIndexed() {
	s.addStoredTx(txtype.Create2Transfer, nil)

	indexValues := s.getToStateIDIndexValues(models.StoredTxName)
	s.Len(indexValues, 1)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
}

// This test checks an edge case that we introduced by indexing ToStateID field which is only available in Transfer transactions.
// See: NewTransactionStorage
func (s *StoredTransactionTestSuite) TestStoredTx_ToStateID_FindUsingIndexWorksWhenThereAreOnlyValuesWithoutThisField() {
	s.addStoredTx(txtype.Create2Transfer, nil)

	txs := make([]models.StoredTx, 0, 1)
	err := s.storage.database.Badger.Find(
		&txs,
		bh.Where("ToStateID").Eq(uint32(1)).Index("ToStateID"),
	)
	s.NoError(err)
	s.Len(txs, 0)
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
	err := s.storage.addStoredTxReceipt(&models.StoredTxReceipt{
		Hash:         utils.RandomHash(),
		CommitmentID: nil, // nil values are not indexed
	})
	s.NoError(err)

	receipts := make([]models.StoredTxReceipt, 0, 1)
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

	indexValues := s.getToStateIDIndexValues(models.StoredTxReceiptName)
	s.Len(indexValues, 3)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[1], 2)
	s.Len(indexValues[2], 1)
}

func (s *StoredTransactionTestSuite) TestStoredTxReceipt_ToStateID_ValuesWithThisFieldSetToNilAreNotIndexed() {
	s.addStoredTxReceipt(nil, nil)

	indexValues := s.getToStateIDIndexValues(models.StoredTxReceiptName)
	s.Len(indexValues, 1)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
}

// This test checks an edge case that we introduced by indexing ToStateID field which can be nil.
// See: NewTransactionStorage
func (s *StoredTransactionTestSuite) TestStoredTxReceipt_ToStateID_FindUsingIndexWorksWhenThereAreOnlyValuesWithThisFieldSetToNil() {
	err := s.storage.addStoredTxReceipt(&models.StoredTxReceipt{
		Hash:      utils.RandomHash(),
		ToStateID: nil, // nil values are not indexed
	})
	s.NoError(err)

	receipts := make([]models.StoredTxReceipt, 0, 1)
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

func (s *StoredTransactionTestSuite) addStoredTx(txType txtype.TransactionType, toStateID *uint32) {
	switch txType {
	case txtype.Transfer:
		err := s.storage.addStoredTx(models.NewStoredTxFromTransfer(&models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash: utils.RandomHash(),
			},
			ToStateID: *toStateID,
		}))
		s.NoError(err)
	case txtype.Create2Transfer:
		err := s.storage.addStoredTx(models.NewStoredTxFromCreate2Transfer(&models.Create2Transfer{
			TransactionBase: models.TransactionBase{
				Hash: utils.RandomHash(),
			},
		}))
		s.NoError(err)
	case txtype.MassMigration:
		panic("not implemented")
	}
}

func (s *StoredTransactionTestSuite) addStoredTxReceipt(toStateID *uint32, commitmentID *models.CommitmentID) {
	receipt := &models.StoredTxReceipt{
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

	s.iterateIndex(models.StoredTxReceiptName, "CommitmentID", func(encodedKey []byte, keyList bh.KeyList) {
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
	indexPrefix := db.IndexKeyPrefix(typeName, indexName)
	err := s.storage.database.Badger.Iterator(indexPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (finish bool, err error) {
		// Get key value
		encodedKeyValue := item.Key()[len(indexPrefix):]

		// Decode value
		var keyList bh.KeyList
		err = item.Value(func(val []byte) error {
			return db.Decode(val, &keyList)
		})
		s.NoError(err)

		handleIndex(encodedKeyValue, keyList)
		return false, nil
	})
	s.ErrorIs(err, db.ErrIteratorFinished)
}

func TestStoredTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(StoredTransactionTestSuite))
}
