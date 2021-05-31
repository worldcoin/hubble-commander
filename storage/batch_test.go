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
		Number:            models.NewUint256(1),
		FinalisationBlock: ref.Uint32(1234),
	}
	id, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(*id)
	s.NoError(err)

	batch.ID = *id
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByNumber() {
	batch := &models.Batch{
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.NewUint256(1234),
		FinalisationBlock: ref.Uint32(1234),
	}
	id, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByNumber(*batch.Number)
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
		Number:            models.NewUint256(1),
		FinalisationBlock: ref.Uint32(1234),
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
		AccountTreeRoot:   nil,
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
		AccountTreeRoot:   nil,
		IncludedInBatch:   nil,
	}

	commitmentID, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	batch, err := s.storage.GetBatchByCommitmentID(*commitmentID)
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(batch)
}

func (s *BatchTestSuite) TestGetOldestPendingBatch() {
	pendingBatches := []models.PendingBatch{
		{
			ID:              2,
			Type:            txtype.Transfer,
			TransactionHash: utils.RandomHash(),
		},
		{
			ID:              3,
			Type:            txtype.Create2Transfer,
			TransactionHash: utils.RandomHash(),
		},
	}
	batch := models.Batch{
		ID:                1,
		Hash:              utils.NewRandomHash(),
		Type:              txtype.Transfer,
		Number:            models.NewUint256(1234),
		FinalisationBlock: ref.Uint32(1234),
	}
	_, err := s.storage.AddBatch(&batch)
	s.NoError(err)
	_, err = s.storage.AddPendingBatch(&pendingBatches[0])
	s.NoError(err)
	_, err = s.storage.AddPendingBatch(&pendingBatches[1])
	s.NoError(err)

	actual, err := s.storage.GetOldestPendingBatch()
	s.NoError(err)

	s.Equal(pendingBatches[0], *actual)
}

func (s *BatchTestSuite) TestGetOldestPendingBatch_NoBatches() {
	res, err := s.storage.GetOldestPendingBatch()
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetLatestSubmittedBatch() {
	batches := []models.Batch{
		{
			ID:                1,
			Hash:              utils.NewRandomHash(),
			Type:              txtype.Transfer,
			Number:            models.NewUint256(1234),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                2,
			Hash:              utils.NewRandomHash(),
			Type:              txtype.Create2Transfer,
			Number:            models.NewUint256(2000),
			FinalisationBlock: ref.Uint32(1234),
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
			Number:            models.NewUint256(1234),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                2,
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			Number:            models.NewUint256(1800),
			FinalisationBlock: ref.Uint32(1800),
		},
		{
			ID:                3,
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			Number:            models.NewUint256(2000),
			FinalisationBlock: ref.Uint32(2000),
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
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.NewUint256(11)},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.NewUint256(12)},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.NewUint256(13)},
		{ID: 4, Hash: utils.NewRandomHash(), Number: models.NewUint256(14)},
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
	_, err := s.storage.AddBatch(&models.Batch{Hash: utils.NewRandomHash(), Number: models.NewUint256(1)})
	s.NoError(err)

	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), models.NewUint256(3))
	s.NoError(err)
	s.Len(actual, 0)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesStartingWithLowerBound() {
	batches := []models.Batch{
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.NewUint256(1)},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.NewUint256(2)},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.NewUint256(3)},
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
		{ID: 1, Hash: utils.NewRandomHash(), Number: models.NewUint256(1)},
		{ID: 2, Hash: utils.NewRandomHash(), Number: models.NewUint256(2)},
		{ID: 3, Hash: utils.NewRandomHash(), Number: models.NewUint256(3)},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(nil, models.NewUint256(2))
	s.NoError(err)
	s.Equal(batches[:2], actual)
}

func (s *BatchTestSuite) TestGetBatchWithAccountRoot_AddAndRetrieve() {
	batch := &models.Batch{
		ID:                1,
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.NewUint256(1),
		FinalisationBlock: ref.Uint32(1234),
	}
	_, err := s.storage.AddBatch(batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.AccountTreeRoot = utils.NewRandomHash()
	includedCommitment.IncludedInBatch = &batch.ID
	_, err = s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	batchWithAccountRoot := &models.BatchWithAccountRoot{
		Batch:           *batch,
		AccountTreeRoot: includedCommitment.AccountTreeRoot,
	}

	actual, err := s.storage.GetBatchWithAccountRoot(*batch.Hash)
	s.NoError(err)
	s.Equal(batchWithAccountRoot, actual)

	actual, err = s.storage.GetBatchWithAccountRootByNumber(*batch.Number)
	s.NoError(err)
	s.Equal(batchWithAccountRoot, actual)
}

func (s *BatchTestSuite) TestGetBatchWithAccountRoot_NotExistingBatch() {
	notFoundErr := NewNotFoundError("batch")
	_, err := s.storage.GetBatchWithAccountRoot(utils.RandomHash())
	s.Equal(notFoundErr, err)

	_, err = s.storage.GetBatchWithAccountRootByNumber(models.MakeUint256(12))
	s.Equal(notFoundErr, err)
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
