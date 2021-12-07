package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	bh "github.com/timshannon/badgerhold/v4"
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
	submissionTime := &models.Timestamp{Time: time.Unix(140, 0).UTC()}
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
		PrevStateRoot:     utils.NewRandomHash(),
		SubmissionTime:    submissionTime,
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(batch.ID)
	s.NoError(err)
	actualUnixTime := actual.SubmissionTime.Unix()

	s.Equal(batch, actual)
	s.EqualValues(140, actualUnixTime)
}

func (s *BatchTestSuite) TestUpdateBatch() {
	pendingBatch := &models.Batch{
		ID:              models.MakeUint256(124),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
	}
	err := s.storage.AddBatch(pendingBatch)
	s.NoError(err)

	batch := &models.Batch{
		ID:                pendingBatch.ID,
		Type:              pendingBatch.Type,
		TransactionHash:   pendingBatch.TransactionHash,
		Hash:              utils.NewRandomHash(),
		SubmissionTime:    &models.Timestamp{Time: time.Unix(140, 0).UTC()},
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err = s.storage.UpdateBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(pendingBatch.ID)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestUpdateBatch_NonexistentBatch() {
	batch := &models.Batch{
		ID:              models.MakeUint256(124),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
	}
	err := s.storage.UpdateBatch(batch)
	s.ErrorIs(err, NewNotFoundError("batch"))
}

func (s *BatchTestSuite) TestGetBatch() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1234),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(batch.ID)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatch_NonexistentBatch() {
	res, err := s.storage.GetBatch(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("batch"))
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetMindedBatch() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1234),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetMinedBatch(batch.ID)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetMinedBatch_PendingBatch() {
	batch := &models.Batch{
		ID:              models.MakeUint256(1234),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		AccountTreeRoot: utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	res, err := s.storage.GetMinedBatch(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("batch"))
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetMinedBatch_NonexistentBatch() {
	res, err := s.storage.GetMinedBatch(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("batch"))
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetLatestSubmittedBatch() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(1234),
			Hash:              utils.NewRandomHash(),
			Type:              batchtype.Transfer,
			FinalisationBlock: ref.Uint32(1234),
			TransactionHash:   utils.RandomHash(),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                models.MakeUint256(2000),
			Hash:              utils.NewRandomHash(),
			Type:              batchtype.Create2Transfer,
			FinalisationBlock: ref.Uint32(1234),
			TransactionHash:   utils.RandomHash(),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
	}
	pendingBatch := models.Batch{
		ID:   models.MakeUint256(2005),
		Type: batchtype.Create2Transfer,
	}
	err := s.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.storage.AddBatch(&batches[1])
	s.NoError(err)
	err = s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	actual, err := s.storage.GetLatestSubmittedBatch()
	s.NoError(err)

	s.Equal(batches[1], *actual)
}

func (s *BatchTestSuite) TestGetLatestSubmittedBatch_NoBatches() {
	res, err := s.storage.GetLatestSubmittedBatch()
	s.ErrorIs(err, NewNotFoundError("batch"))
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetLatestFinalisedBatch() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(1234),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                models.MakeUint256(1800),
			Type:              batchtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1800),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
		{
			ID:                models.MakeUint256(2000),
			Type:              batchtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2000),
			AccountTreeRoot:   utils.NewRandomHash(),
		},
	}
	pendingBatch := models.Batch{
		ID:              models.MakeUint256(2005),
		Type:            batchtype.Create2Transfer,
		TransactionHash: utils.RandomHash(),
	}

	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}

	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	finalisedBatch, err := s.storage.GetLatestFinalisedBatch(1800)
	s.NoError(err)

	s.Equal(batches[1], *finalisedBatch)
}

func (s *BatchTestSuite) TestGetLatestFinalisedBatch_NoBatches() {
	res, err := s.storage.GetLatestFinalisedBatch(500)
	s.ErrorIs(err, NewNotFoundError("batch"))
	s.Nil(res)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsCorrectBatches() {
	batches := []models.Batch{
		{ID: models.MakeUint256(11), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(12), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(13), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(14), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(models.NewUint256(12), models.NewUint256(13))
	s.NoError(err)
	s.Equal(batches[1:3], actual)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsEmptySliceWhenThereAreNoBatchesInRange() {
	err := s.storage.AddBatch(&models.Batch{ID: models.MakeUint256(1), Hash: utils.NewRandomHash()})
	s.NoError(err)

	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), models.NewUint256(3))
	s.NoError(err)
	s.Len(actual, 0)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesStartingWithLowerBound() {
	batches := []models.Batch{
		{ID: models.MakeUint256(1), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(2), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(3), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(models.NewUint256(2), nil)
	s.NoError(err)
	s.Equal(batches[1:], actual)
}

func (s *BatchTestSuite) TestGetBatchesInRange_ReturnsAllBatchesUpUntilUpperBound() {
	batches := []models.Batch{
		{ID: models.MakeUint256(1), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(2), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
		{ID: models.MakeUint256(3), Hash: utils.NewRandomHash(), TransactionHash: utils.RandomHash()},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}
	actual, err := s.storage.GetBatchesInRange(nil, models.NewUint256(2))
	s.NoError(err)
	s.Equal(batches[:2], actual)
}

func (s *BatchTestSuite) TestGetBatchByHash_AddAndRetrieve() {
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   utils.NewRandomHash(),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByHash(*batch.Hash)
	s.NoError(err)
	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByHash_NonexistentBatch() {
	_, err := s.storage.GetBatchByHash(utils.RandomHash())
	s.True(IsNotFoundError(err))
}

func (s *BatchTestSuite) TestDeleteBatches() {
	batches := []models.Batch{
		{
			ID:              models.MakeUint256(1),
			Type:            batchtype.Transfer,
			TransactionHash: utils.RandomHash(),
			Hash:            utils.NewRandomHash(),
		},
		{
			ID:              models.MakeUint256(2),
			Type:            batchtype.Create2Transfer,
			TransactionHash: utils.RandomHash(),
			Hash:            utils.NewRandomHash(),
		},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)
	}

	err := s.storage.DeleteBatches(batches[0].ID, batches[1].ID)
	s.NoError(err)

	for i := range batches {
		_, err = s.storage.GetBatch(batches[i].ID)
		s.ErrorIs(err, NewNotFoundError("batch"))
	}
}

func (s *BatchTestSuite) TestDeleteBatches_NotExistentBatch() {
	err := s.storage.DeleteBatches(models.MakeUint256(1))
	s.ErrorIs(err, NewNotFoundError("batch"))
}

func (s *BatchTestSuite) TestBatch_Hash_IndexWorks() {
	hash1 := utils.NewRandomHash()
	hash2 := utils.NewRandomHash()
	s.addBatch(1, hash1)
	s.addBatch(2, hash2)
	s.addBatch(3, hash1)

	indexValues := s.getHashIndexValues()
	s.Len(indexValues, 3)
	s.Len(indexValues[common.Hash{}], 0) // value set due to index initialization, see NewTransactionStorage
	s.Len(indexValues[*hash1], 2)
	s.Len(indexValues[*hash2], 1)
}

func (s *BatchTestSuite) TestBatch_Hash_ValuesWithThisFieldSetToNilAreNotIndexed() {
	s.addBatch(1, nil)

	indexValues := s.getHashIndexValues()
	s.Len(indexValues, 1)
	s.Len(indexValues[common.Hash{}], 0) // value set due to index initialization, see NewTransactionStorage
}

func (s *BatchTestSuite) TestBatch_Hash_FindUsingIndexWorksWhenThereAreOnlyValuesWithThisFieldSetToNil() {
	s.addBatch(1, nil)

	txs := make([]models.Batch, 0, 1)
	err := s.storage.database.Badger.Find(
		&txs,
		bh.Where("Hash").Eq(utils.RandomHash()).Index("Hash"),
	)
	s.NoError(err)
	s.Len(txs, 0)
}

func (s *BatchTestSuite) getHashIndexValues() map[common.Hash]bh.KeyList {
	indexValues := make(map[common.Hash]bh.KeyList)

	s.iterateIndex(stored.BatchName, "Hash", func(encodedKey []byte, keyList bh.KeyList) {
		var batchHash common.Hash
		err := db.Decode(encodedKey, &batchHash)
		s.NoError(err)

		indexValues[batchHash] = keyList
	})

	return indexValues
}

func (s *BatchTestSuite) addBatch(id uint64, hash *common.Hash) {
	err := s.storage.AddBatch(&models.Batch{
		ID:   models.MakeUint256(id),
		Hash: hash,
	})
	s.NoError(err)
}

func (s *BatchTestSuite) iterateIndex(
	typeName []byte,
	indexName string,
	handleIndex func(encodedKey []byte, keyList bh.KeyList),
) {
	testutils.IterateIndex(s.Assertions, s.storage.database.Badger, typeName, indexName, handleIndex)
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
