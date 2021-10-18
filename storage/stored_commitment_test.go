package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoredCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StoredCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StoredCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment() {
	expected := commitment
	for i := 0; i < 2; i++ {
		expected.ID.IndexInBatch = uint8(i)
		err := s.storage.AddCommitment(&expected)
		s.NoError(err)
	}

	expected.ID.BatchID = models.MakeUint256(5)
	for i := 0; i < 2; i++ {
		expected.ID.IndexInBatch = uint8(i)
		err := s.storage.AddCommitment(&expected)
		s.NoError(err)
	}

	latestCommitment, err := s.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(models.MakeStoredCommitmentFromTxCommitment(&expected), *latestCommitment)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment_NoCommitments() {
	_, err := s.storage.GetLatestCommitment()
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func (s *StoredCommitmentTestSuite) TestDeleteCommitmentsByBatchIDs() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(111),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                models.MakeUint256(5),
			Type:              batchtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2345),
		},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)

		for j := 0; j < 2; j++ {
			commitmentInBatch := commitment
			commitmentInBatch.ID.BatchID = batches[i].ID
			commitmentInBatch.ID.IndexInBatch = uint8(j)
			err = s.storage.AddCommitment(&commitmentInBatch)
			s.NoError(err)
		}
	}

	err := s.storage.DeleteCommitmentsByBatchIDs(batches[0].ID, batches[1].ID)
	s.NoError(err)
	for i := range batches {
		_, err = s.storage.getStoredCommitmentsByBatchID(batches[i].ID)
		s.ErrorIs(err, NewNotFoundError("commitments"))
	}
}

func (s *StoredCommitmentTestSuite) TestDeleteCommitmentsByBatchIDs_NoCommitments() {
	batchID := s.addRandomBatch()
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.DeleteCommitmentsByBatchIDs(batchID)
	s.ErrorIs(err, NewNotFoundError("commitments"))

	_, err = s.storage.GetCommitment(&commitment.ID)
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func TestStoredCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(StoredCommitmentTestSuite))
}
