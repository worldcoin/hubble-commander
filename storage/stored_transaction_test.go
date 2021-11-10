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

func (s *StoredTransactionTestSuite) TestSetTransactionError() {
	err := s.storage.AddTransfer(&transferTransaction)
	s.NoError(err)

	errorMessage := ref.String("Quack")

	err = s.storage.SetTransactionError(transferTransaction.Hash, *errorMessage)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transferTransaction.Hash)
	s.NoError(err)

	s.Equal(errorMessage, res.ErrorMessage)
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

	err := s.storage.AddTransfer(&tx1)
	s.NoError(err)
	err = s.storage.AddCreate2Transfer(&tx2)
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

	c2t := create2Transfer
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.CommitmentID = &commitmentInBatch.ID
	err = s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	count, err := s.storage.GetTransactionCount()
	s.NoError(err)
	s.Equal(2, *count)
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

func (s *StoredTransactionTestSuite) TestAddStoredTxReceipt_IndexOnToStateIDWorks() {
	s.addStoredTxReceipt(ref.Uint32(1))
	s.addStoredTxReceipt(ref.Uint32(2))
	s.addStoredTxReceipt(ref.Uint32(1))

	indexValues := s.getToStateIDIndexValues()
	s.Len(indexValues, 3)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[1], 2)
	s.Len(indexValues[2], 1)
}

func (s *StoredTransactionTestSuite) TestAddStoredTxReceipt_ValuesWithNilToStateIDAreNotIndexed() {
	s.addStoredTxReceipt(nil)

	indexValues := s.getToStateIDIndexValues()
	s.Len(indexValues, 1)
	s.Len(indexValues[0], 0) // value set due to index initialization, see NewTransactionStorage
}

// TODO do the same test for StoredTx
// This test checks an edge case that we introduced by indexing ToStateID field which can be nil.
// See: NewTransactionStorage
func (s *StoredTransactionTestSuite) TestStoredTxReceipt_FindUsingIndexWorksWhenThereAreOnlyStoredTxReceiptsWithNilToStateID() {
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

func (s *StoredTransactionTestSuite) addStoredTxReceipt(toStateID *uint32) {
	receipt := &models.StoredTxReceipt{
		Hash:      utils.RandomHash(),
		ToStateID: toStateID,
	}
	err := s.storage.addStoredTxReceipt(receipt)
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) getToStateIDIndexValues() map[uint32]bh.KeyList {
	indexValues := make(map[uint32]bh.KeyList)

	indexPrefix := db.IndexKeyPrefix(models.StoredTxReceiptName, "ToStateID")
	err := s.storage.database.Badger.Iterator(indexPrefix, db.PrefetchIteratorOpts, func(item *bdg.Item) (finish bool, err error) {
		// Decode key
		encodedToStateID := item.Key()[len(indexPrefix):]
		var toStateID uint32
		err = db.Decode(encodedToStateID, &toStateID)
		s.NoError(err)

		// Decode value
		var keyList bh.KeyList
		err = item.Value(func(val []byte) error {
			return db.Decode(val, &keyList)
		})
		s.NoError(err)

		indexValues[toStateID] = keyList
		return false, nil
	})
	s.ErrorIs(err, db.ErrIteratorFinished)

	return indexValues
}

func TestStoredTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(StoredTransactionTestSuite))
}
