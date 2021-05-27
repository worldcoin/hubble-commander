package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
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
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
		Number:            models.MakeUint256(1),
		FinalisationBlock: 1234,
	}
	_, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(batch.Hash)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByID() {
	batch := &models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
		Number:            models.MakeUint256(1234),
		FinalisationBlock: 1234,
	}
	_, err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByID(batch.Number)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatch_NonExistentBatch() {
	res, err := s.storage.GetBatch(common.Hash{1, 2, 3, 4})
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatchByID_NonExistentBatch() {
	res, err := s.storage.GetBatchByID(models.MakeUint256(42))
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatchByCommitmentID() {
	batchHash := utils.RandomHash()

	batch := &models.Batch{
		Hash:              batchHash,
		Type:              txtype.Transfer,
		Number:            models.MakeUint256(1),
		FinalisationBlock: 1234,
	}

	_, err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitment := &models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		AccountTreeRoot:   nil,
		IncludedInBatch:   &batchHash,
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

func (s *BatchTestSuite) TestGetLatestBatch() {
	batches := []models.Batch{
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Transfer,
			Number:            models.MakeUint256(1234),
			FinalisationBlock: 1234,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			Number:            models.MakeUint256(2000),
			FinalisationBlock: 2000,
		},
	}
	_, err := s.storage.AddBatch(&batches[0])
	s.NoError(err)
	_, err = s.storage.AddBatch(&batches[1])
	s.NoError(err)

	actual, err := s.storage.GetLatestBatch()
	s.NoError(err)

	s.Equal(batches[1], *actual)
}

func (s *BatchTestSuite) TestGetLatestBatch_NoBatches() {
	res, err := s.storage.GetLatestBatch()
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetLatestFinalisedBatch() {
	batches := []models.Batch{
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Transfer,
			Number:            models.MakeUint256(1234),
			FinalisationBlock: 1234,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			Number:            models.MakeUint256(1800),
			FinalisationBlock: 1800,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			Number:            models.MakeUint256(2000),
			FinalisationBlock: 2000,
		},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}

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
		{Hash: utils.RandomHash(), Number: models.MakeUint256(1)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(2)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(3)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(4)},
	}
	for i := range batches {
		_, err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), models.NewUint256(3))
	s.NoError(err)
	s.Equal(batches[1:3], actual)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsEmptySliceWhenThereAreNoBatchesInRange() {
	_, err := s.storage.AddBatch(&models.Batch{Hash: utils.RandomHash(), Number: models.MakeUint256(1)})
	s.NoError(err)

	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), models.NewUint256(3))
	s.NoError(err)
	s.Len(actual, 0)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesStartingWithLowerBound() {
	batches := []models.Batch{
		{Hash: utils.RandomHash(), Number: models.MakeUint256(1)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(2)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(3)},
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
		{Hash: utils.RandomHash(), Number: models.MakeUint256(1)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(2)},
		{Hash: utils.RandomHash(), Number: models.MakeUint256(3)},
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
	hash := utils.RandomHash()
	batch := &models.BatchWithAccountRoot{
		BatchWithSubmissionBlock: models.BatchWithSubmissionBlock{
			Batch: models.Batch{
				Hash:              hash,
				Type:              txtype.Transfer,
				Number:            models.MakeUint256(1),
				FinalisationBlock: 1234,
			},
		},
		AccountTreeRoot: &hash,
	}
	_, err := s.storage.AddBatch(&batch.Batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.AccountTreeRoot = &hash
	includedCommitment.IncludedInBatch = &batch.Hash
	_, err = s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	actual, err := s.storage.GetBatchWithAccountRoot(batch.Hash)
	s.NoError(err)
	s.Equal(batch, actual)

	actual, err = s.storage.GetBatchWithAccountRootByID(batch.Number)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchWithAccountRoot_NotExistingBatch() {
	notFoundErr := NewNotFoundError("batch")
	_, err := s.storage.GetBatchWithAccountRoot(utils.RandomHash())
	s.Equal(notFoundErr, err)

	_, err = s.storage.GetBatchWithAccountRootByID(models.MakeUint256(12))
	s.Equal(notFoundErr, err)
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
