package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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
		SubmissionTime:  &models.Timestamp{Time: time.Unix(140, 0).UTC()},
	}

	err = s.storage.AddBatch(s.batch)
	s.NoError(err)
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) TestGetTransactionWithBatchDetails_WithoutBatch() {
	err := s.storage.AddTransfer(&transfer)
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
	err := s.storage.AddTransfer(&transferInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: transferInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		BatchTime:   s.batch.SubmissionTime,
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
	err := s.storage.AddCreate2Transfer(&create2TransferInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: create2TransferInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		BatchTime:   s.batch.SubmissionTime,
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
	err := s.storage.AddMassMigration(&massMigrationInBatch)
	s.NoError(err)

	expected := models.TransactionWithBatchDetails{
		Transaction: massMigrationInBatch.Copy(),
		BatchHash:   s.batch.Hash,
		BatchTime:   s.batch.SubmissionTime,
	}
	res, err := s.storage.GetTransactionWithBatchDetails(massMigrationInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
