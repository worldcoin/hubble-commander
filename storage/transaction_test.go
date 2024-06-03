package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
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
		PrevStateRoot:   utils.RandomHash(),
	}

	err = s.storage.AddBatch(s.batch)
	s.NoError(err)

	leaf := models.StateLeaf{
		StateID: 0,
		UserState: models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(2000),
		},
	}
	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	leaf = models.StateLeaf{
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

func (s *TransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_WithoutBatch() {
	transfer = models.Transfer{
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

	err := s.storage.AddMempoolTx(&transfer)
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
	transferInBatch.CommitmentSlot = &models.CommitmentSlot{
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
	create2TransferInBatch.CommitmentSlot = &models.CommitmentSlot{
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
	massMigrationInBatch.CommitmentSlot = &models.CommitmentSlot{
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

func (s *TransactionTestSuite) TestGetTransactionsByCommitmentID() {
	transfer1 := transfer
	transfer1.CommitmentSlot = models.NewCommitmentSlot(txCommitment.ID, 0)
	err := s.storage.AddTransaction(&transfer1)
	s.NoError(err)

	otherCommitmentID := txCommitment.ID
	otherCommitmentID.IndexInBatch += 1
	transfer2 := create2Transfer
	transfer2.Hash = utils.RandomHash()
	transfer2.CommitmentSlot = models.NewCommitmentSlot(otherCommitmentID, 0)
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

func (s *TransactionTestSuite) TestBatchAdd_Nothing() {
	err := s.storage.BatchAddTransaction(models.MakeMassMigrationArray())
	s.ErrorIs(err, ErrNoRowsAffected)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
