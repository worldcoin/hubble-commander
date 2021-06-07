package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BatchTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *BatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *BatchTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *BatchTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *BatchTestSuite) TestAddBatch_AddAndRetrieve() {
	batch := &models.Batch{
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	id, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(*id)
	s.NoError(err)

	batch.ID = *id
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestMarkBatchAsSubmitted() {
	pendingBatch := &models.Batch{
		Type:            txtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Number:          models.MakeUint256(124),
	}
	batchID, err := s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	batch := &models.Batch{
		ID:                *batchID,
		Type:              pendingBatch.Type,
		TransactionHash:   pendingBatch.TransactionHash,
		Hash:              utils.NewRandomHash(),
		Number:            pendingBatch.Number,
		FinalisationBlock: ref.Uint32(1234),
	}
	err = s.storage.MarkBatchAsSubmitted(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(*batchID)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByNumber() {
	batch := &models.Batch{
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.MakeUint256(1234),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	id, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByNumber(batch.Number)
	s.NoError(err)

	batch.ID = *id
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByNumber_NonExistentBatch() {
	res, err := s.storage.GetBatchByNumber(models.MakeUint256(42))
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatch_NonExistentBatch() {
	res, err := s.storage.GetBatch(1)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatchByCommitmentID() {
	batch := &models.Batch{
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}

	batchID, err := s.storage.AddBatch(batch)
	s.NoError(err)

	batch.ID = *batchID

	commitment := &models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		IncludedInBatch:   &batch.ID,
	}

	commitmentID, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	actual, err := s.storage.GetBatchByCommitmentID(*commitmentID)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByCommitmentID_NotExistentBatch() {
	commitment := &models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		IncludedInBatch:   nil,
	}

	commitmentID, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	batch, err := s.storage.GetBatchByCommitmentID(*commitmentID)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(batch)
}

func (s *BatchTestSuite) TestGetLatestSubmittedBatch() {
	batches := []models.Batch{
		{
			ID:                1,
			Hash:              utils.NewRandomHash(),
			Type:              txtype.Transfer,
			Number:            models.MakeUint256(1234),
			FinalisationBlock: ref.Uint32(1234),
			TransactionHash:   utils.RandomHash(),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                2,
			Hash:              utils.NewRandomHash(),
			Type:              txtype.Create2Transfer,
			Number:            models.MakeUint256(2000),
			FinalisationBlock: ref.Uint32(1234),
			TransactionHash:   utils.RandomHash(),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
	}
	pendingBatch := models.Batch{
		ID:   3,
		Type: txtype.Create2Transfer,
	}
	_, err := s.storage.AddBatch(&batches[0])
	s.NoError(err)
	_, err = s.storage.AddBatch(&batches[1])
	s.NoError(err)
	_, err = s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	actual, err := s.storage.GetLatestSubmittedBatch()
	s.NoError(err)

	s.Equal(batches[1], *actual)
}

func (s *BatchTestSuite) TestGetLatestSubmittedBatch_NoBatches() {
	res, err := s.storage.GetLatestSubmittedBatch()
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetLatestFinalisedBatch() {
	batches := []models.Batch{
		{
			ID:                1,
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			Number:            models.MakeUint256(1234),
			FinalisationBlock: ref.Uint32(1234),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                2,
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			Number:            models.MakeUint256(1800),
			FinalisationBlock: ref.Uint32(1800),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                3,
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			Number:            models.MakeUint256(2000),
			FinalisationBlock: ref.Uint32(2000),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
	}
	pendingBatch := models.Batch{
		ID:              4,
		Type:            txtype.Create2Transfer,
		TransactionHash: utils.RandomHash(),
	}

	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}

	_, err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	finalisedBatch, err := s.storage.GetLatestFinalisedBatch(1800)
	s.NoError(err)

	s.Equal(batches[1], *finalisedBatch)
}

func (s *BatchTestSuite) TestGetLatestFinalisedBatch_NoBatches() {
	res, err := s.storage.GetLatestFinalisedBatch(500)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsCorrectBatches() {
	batches := []models.Batch{
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.MakeUint256(11), TransactionHash: utils.RandomHash()},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.MakeUint256(12), TransactionHash: utils.RandomHash()},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.MakeUint256(13), TransactionHash: utils.RandomHash()},
		{ID: 4, Hash: utils.NewRandomHash(), Number: models.MakeUint256(14), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(models.NewUint256(12), models.NewUint256(13))
	s.NoError(err)
	s.Equal(batches[1:3], actual)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsEmptySliceWhenThereAreNoBatchesInRange() {
	_, err := s.storage.AddBatch(&models.Batch{Hash: utils.NewRandomHash(), Number: models.MakeUint256(1)})
	s.NoError(err)

	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), models.NewUint256(3))
	s.NoError(err)
	s.Len(actual, 0)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesStartingWithLowerBound() {
	batches := []models.Batch{
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.MakeUint256(1), TransactionHash: utils.RandomHash()},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.MakeUint256(2), TransactionHash: utils.RandomHash()},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.MakeUint256(3), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), nil)
	s.NoError(err)
	s.Equal(batches[1:], actual)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesUpUntilUpperBound() {
	batches := []models.Batch{
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.MakeUint256(1), TransactionHash: utils.RandomHash()},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.MakeUint256(2), TransactionHash: utils.RandomHash()},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.MakeUint256(3), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(nil, models.NewUint256(2))
	s.NoError(err)
	s.Equal(batches[:2], actual)
}

func (s *BatchTestSuite) TestGetBatchByHash_AddAndRetrieve() {
	batch := &models.Batch{
		ID:                1,
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	_, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByHash(*batch.Hash)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByHash_NotExistingBatch() {
	_, err := s.storage.GetBatchByHash(utils.RandomHash())
	s.True(IsNotFoundError(err))
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
