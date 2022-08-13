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
	storage           *TestStorage
	depositCommitment models.DepositCommitment
}

func (s *StoredCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.depositCommitment = models.DepositCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(5),
				IndexInBatch: 1,
			},
			Type: batchtype.Deposit,
		},
		SubtreeID:   models.MakeUint256(1),
		SubtreeRoot: utils.RandomHash(),
		Deposits: []models.PendingDeposit{
			{
				ID: models.DepositID{
					SubtreeID:    models.MakeUint256(1),
					DepositIndex: models.MakeUint256(0),
				},
			},
		},
	}
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
	expected := txCommitment
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
	s.Equal(expected.CommitmentBase, *latestCommitment)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment_LatestDepositCommitment() {
	for i := 0; i < 2; i++ {
		commitment := txCommitment
		commitment.ID.IndexInBatch = uint8(i)
		err := s.storage.AddCommitment(&commitment)
		s.NoError(err)
	}

	err := s.storage.AddCommitment(&s.depositCommitment)
	s.NoError(err)

	latestCommitment, err := s.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(s.depositCommitment.CommitmentBase, *latestCommitment)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment_NoCommitments() {
	_, err := s.storage.GetLatestCommitment()
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func (s *StoredCommitmentTestSuite) TestRemoveCommitmentsByBatchIDs() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(111),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
			PrevStateRoot:     utils.RandomHash(),
		},
		{
			ID:                models.MakeUint256(5),
			Type:              batchtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2345),
			PrevStateRoot:     utils.RandomHash(),
		},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)

		for j := 0; j < 2; j++ {
			commitmentInBatch := txCommitment
			commitmentInBatch.ID.BatchID = batches[i].ID
			commitmentInBatch.ID.IndexInBatch = uint8(j)
			err = s.storage.AddCommitment(&commitmentInBatch)
			s.NoError(err)
		}
	}

	depositCommitment := s.depositCommitment
	depositCommitment.ID = models.CommitmentID{
		BatchID:      batches[0].ID,
		IndexInBatch: 2,
	}
	err := s.storage.AddCommitment(&depositCommitment)
	s.NoError(err)

	err = s.storage.RemoveCommitmentsByBatchIDs(batches[0].ID, batches[1].ID)
	s.NoError(err)
	for i := range batches {
		commitments, err := s.storage.getStoredCommitmentsByBatchID(batches[i].ID)
		s.NoError(err)
		s.Len(commitments, 0)
	}
}

func (s *StoredCommitmentTestSuite) TestRemoveCommitmentsByBatchIDs_NoCommitments() {
	batchID := s.addRandomBatch()
	err := s.storage.AddCommitment(&txCommitment)
	s.NoError(err)

	err = s.storage.RemoveCommitmentsByBatchIDs(batchID)
	s.ErrorIs(err, NewNotFoundError("commitments"))

	_, err = s.storage.GetCommitment(&txCommitment.ID)
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		PrevStateRoot:     utils.RandomHash(),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func TestStoredCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(StoredCommitmentTestSuite))
}
