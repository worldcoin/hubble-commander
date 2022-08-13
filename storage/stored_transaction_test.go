package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
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

	leaf := models.StateLeaf{
		StateID: 1,
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(2000),
		},
	}
	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	leaf = models.StateLeaf{
		StateID: 2,
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(2000),
		},
	}
	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredTransactionTestSuite) TestMarkTransactionsAsPending() {
	// this test fails and it should fail because the current behavior is a bug,
	// transactions should be put back into the mempool when a batch reverts
	s.T().Skip()

	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transfer
		txs[i].Hash = utils.RandomHash()
		txs[i].CommitmentSlot = &models.CommitmentSlot{
			BatchID:           models.MakeUint256(5),
			IndexInBatch:      3,
			IndexInCommitment: uint8(i),
		}
		err := s.storage.AddTransaction(&txs[i])
		s.NoError(err)
	}

	err := s.storage.MarkTransactionsAsPending([]models.CommitmentSlot{
		*txs[0].CommitmentSlot, *txs[1].CommitmentSlot},
	)
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		s.Nil(tx.CommitmentSlot)
	}
}

func (s *StoredTransactionTestSuite) TestMarkTransactionsAsIncluded() {
	tx1 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: 2,
	}
	err := s.storage.AddMempoolTx(&tx1)
	s.NoError(err)

	tx2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(1),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: 2,
	}
	err = s.storage.AddMempoolTx(&tx2)
	s.NoError(err)

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 1,
	}
	err = s.storage.MarkTransactionsAsIncluded(models.MakeTransferArray(tx1, tx2), &commitmentID)
	s.NoError(err)

	tx, err := s.storage.GetTransfer(tx1.Hash)
	s.NoError(err)
	s.Equal(commitmentID, *tx.CommitmentSlot.CommitmentID())

	tx, err = s.storage.GetTransfer(tx2.Hash)
	s.NoError(err)
	s.Equal(commitmentID, *tx.CommitmentSlot.CommitmentID())
}

func (s *StoredTransactionTestSuite) TestGetTransactionCount() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		PrevStateRoot:     utils.RandomHash(),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := txCommitment
	commitmentInBatch.ID.BatchID = batch.ID
	err = s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInCommitment := transfer
	transferInCommitment.Hash = common.Hash{5, 5, 5}
	transferInCommitment.CommitmentSlot = models.NewCommitmentSlot(commitmentInBatch.ID, 1)
	err = s.storage.AddTransaction(&transferInCommitment)
	s.NoError(err)

	c2t := create2Transfer
	c2t.Hash = common.Hash{3, 4, 5}
	c2t.CommitmentSlot = models.NewCommitmentSlot(commitmentInBatch.ID, 2)
	err = s.storage.AddTransaction(&c2t)
	s.NoError(err)

	mm := massMigration
	mm.Hash = common.Hash{6, 7, 8}
	mm.CommitmentSlot = models.NewCommitmentSlot(commitmentInBatch.ID, 3)
	err = s.storage.AddTransaction(&mm)
	s.NoError(err)

	storageCount, err := s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(3, *storageCount)
	count := s.storage.GetTransactionCount()
	s.EqualValues(3, count)

	err = s.storage.MarkTransactionsAsPending([]models.CommitmentSlot{*transferInCommitment.CommitmentSlot})
	s.NoError(err)
	storageCount, err = s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(2, *storageCount)
	count = s.storage.GetTransactionCount()
	s.EqualValues(2, count)

	err = s.storage.MarkTransactionsAsIncluded(models.MakeTransferArray(transferInCommitment), &commitmentInBatch.ID)
	s.NoError(err)
	storageCount, err = s.storage.getTransactionCount()
	s.NoError(err)
	s.EqualValues(3, *storageCount)
	count = s.storage.GetTransactionCount()
	s.EqualValues(3, count)
}

func (s *StoredTransactionTestSuite) TestGetTransactionCount_IncrementsTxCountOnStorageCopy() {
	transactionStorageCopy := s.storage.TransactionStorage.copyWithNewDatabase(s.storage.database)
	transactionStorageCopy.incrementTransactionCount()

	count := s.storage.GetTransactionCount()
	s.EqualValues(1, count)
}

func (s *StoredTransactionTestSuite) TestGetTransactionCount_NoTransactions() {
	count := s.storage.GetTransactionCount()
	s.EqualValues(0, count)
}

func (s *StoredTransactionTestSuite) TestGetTransactionIDsByBatchIDs() {
	batchIDs := []models.Uint256{models.MakeUint256(1), models.MakeUint256(2)}
	expectedIDs := make([]models.CommitmentSlot, 0, 4)
	for i := range batchIDs {
		transfers := make([]models.Transfer, 2)
		transfers[0] = transfer
		transfers[0].Hash = utils.RandomHash()
		transfers[1] = transfer
		transfers[1].Hash = utils.RandomHash()
		s.addTransfersInCommitment(&batchIDs[i], transfers)
		expectedIDs = append(expectedIDs, *transfers[0].CommitmentSlot, *transfers[1].CommitmentSlot)
	}

	ids, err := s.storage.GetTransactionIDsByBatchIDs(batchIDs...)
	s.NoError(err)
	s.Len(ids, 4)
	for i := range expectedIDs {
		s.Contains(ids, expectedIDs[i])
	}
}

func (s *StoredTransactionTestSuite) TestGetTransactionIDsByBatchIDs_NoTransactions() {
	transfers := make([]models.Transfer, 2)
	transfers[0] = transfer
	transfers[1] = transfer
	transfers[1].Hash = utils.RandomHash()
	s.addTransfersInCommitment(models.NewUint256(1), transfers)

	ids, err := s.storage.GetTransactionIDsByBatchIDs(models.MakeUint256(2))
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(ids)
}

func (s *StoredTransactionTestSuite) TestGetAllPendingTransactions() {
	expectedTxs := s.populatePendingTransactions()

	res, err := s.storage.GetAllMempoolTransactions()
	s.NoError(err)
	s.Len(res, 3)
	for i := range expectedTxs {
		expectedTx := expectedTxs[i]
		s.Contains(res, expectedTx)
	}
}

func (s *StoredTransactionTestSuite) TestGetAllPendingTransactions_NoTransactions() {
	txs, err := s.storage.GetAllMempoolTransactions()
	s.NoError(err)
	s.Len(txs, 0)
}

func First(txs []models.GenericTransaction, pred func(tx models.GenericTransaction) bool) models.GenericTransaction {
	for i := range txs {
		if pred(txs[i]) {
			return txs[i]
		}
	}

	return nil
}

func (s *StoredTransactionTestSuite) TestGetAllFailedTransactions() {
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

	transfer := First(res, func(tx models.GenericTransaction) bool { return tx.Type() == txtype.Transfer })
	s.Equal(expectedTxs[0], transfer.ToTransfer())

	c2t := First(res, func(tx models.GenericTransaction) bool { return tx.Type() == txtype.Create2Transfer })
	s.Equal(expectedTxs[1], c2t.ToCreate2Transfer())

	mm := First(res, func(tx models.GenericTransaction) bool { return tx.Type() == txtype.MassMigration })
	s.Equal(expectedTxs[2], mm.ToMassMigration())
}

func (s *StoredTransactionTestSuite) TestGetAllFailedTransactions_NoTransactions() {
	txs, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *StoredTransactionTestSuite) populatePendingTransactions() []stored.PendingTx {
	transfer1 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: 2,
	}
	err := s.storage.AddMempoolTx(&transfer1)
	s.NoError(err)

	create2Transfer1 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(1),
			Signature:   models.MakeRandomSignature(),
		},
		ToPublicKey: account2.PublicKey,
	}
	err = s.storage.AddMempoolTx(&create2Transfer1)
	s.NoError(err)

	massMigration1 := models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: 1,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(2),
			Signature:   models.MakeRandomSignature(),
		},
		SpokeID: 5,
	}
	err = s.storage.AddMempoolTx(&massMigration1)
	s.NoError(err)

	return []stored.PendingTx{
		*stored.NewPendingTx(&transfer1),
		*stored.NewPendingTx(&create2Transfer1),
		*stored.NewPendingTx(&massMigration1),
	}
}

func (s *StoredTransactionTestSuite) populateTransactions() models.GenericTransactionArray {
	transfer1 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
			CommitmentSlot: &models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
		},
		ToStateID: 2,
	}

	transfers := make([]models.Transfer, 3)
	for i := range transfers {
		transfers[i] = transfer1
		transfers[i].Hash = utils.RandomHash()
	}
	transfers[1].CommitmentSlot = &models.CommitmentSlot{BatchID: models.MakeUint256(2)}
	transfers[2].ErrorMessage = ref.String("A very boring error message")
	transfers[2].CommitmentSlot = nil

	err := s.storage.BatchAddTransaction(models.MakeTransferArray(transfers...))
	s.NoError(err)

	create2Transfer1 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
			CommitmentSlot: &models.CommitmentSlot{
				BatchID:           models.MakeUint256(3),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
		},
		ToPublicKey: account2.PublicKey,
	}

	create2Transfers := make([]models.Create2Transfer, 3)
	for i := range create2Transfers {
		create2Transfers[i] = create2Transfer1
		create2Transfers[i].Hash = utils.RandomHash()
	}

	create2Transfers[1].CommitmentSlot = &models.CommitmentSlot{BatchID: models.MakeUint256(4)}
	create2Transfers[2].ErrorMessage = ref.String("A very boring error message")
	create2Transfers[2].CommitmentSlot = nil

	err = s.storage.BatchAddTransaction(models.MakeCreate2TransferArray(create2Transfers...))
	s.NoError(err)

	massMigration1 := models.MassMigration{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.MassMigration,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
			CommitmentSlot: &models.CommitmentSlot{
				BatchID:           models.MakeUint256(5),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
		},
		SpokeID: 5,
	}

	massMigrations := make([]models.MassMigration, 3)
	for i := range massMigrations {
		massMigrations[i] = massMigration1
		massMigrations[i].Hash = utils.RandomHash()
	}
	massMigrations[1].CommitmentSlot = &models.CommitmentSlot{BatchID: models.MakeUint256(6)}
	massMigrations[2].ErrorMessage = ref.String("A very boring error message")
	massMigrations[2].CommitmentSlot = nil

	err = s.storage.BatchAddTransaction(models.MakeMassMigrationArray(massMigrations...))
	s.NoError(err)

	var result models.GenericTransactionArray
	result = models.MakeGenericArray()
	result = result.Append(models.MakeTransferArray(transfers...))
	result = result.Append(models.MakeCreate2TransferArray(create2Transfers...))
	result = result.Append(models.MakeMassMigrationArray(massMigrations...))
	return result
}

func (s *StoredTransactionTestSuite) addTransfersInCommitment(batchID *models.Uint256, transfers []models.Transfer) {
	for i := range transfers {
		transfers[i].CommitmentSlot = &models.CommitmentSlot{
			BatchID:      *batchID,
			IndexInBatch: uint8(i),
		}
		err := s.storage.AddTransaction(&transfers[i])
		s.NoError(err)
	}
}

func TestStoredTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(StoredTransactionTestSuite))
}
