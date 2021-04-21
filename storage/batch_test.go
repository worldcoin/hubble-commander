package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
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
	storage *Storage
	db      *db.TestDB
}

func (s *BatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *BatchTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *BatchTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *BatchTestSuite) TestAddBatch_AddAndRetrieve() {
	batch := &models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
		ID:                models.MakeUint256(1),
		FinalisationBlock: 1234,
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(batch.Hash)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) TestGetBatchByID() {
	batch := &models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
		ID:                models.MakeUint256(1234),
		FinalisationBlock: 1234,
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByID(batch.ID)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *StateUpdateTestSuite) TestGetBatch_NonExistentBatch() {
	res, err := s.storage.GetBatch(common.Hash{1, 2, 3, 4})
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *StateUpdateTestSuite) TestGetBatchByID_NonExistentBatch() {
	res, err := s.storage.GetBatchByID(models.MakeUint256(42))
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}

func (s *StateUpdateTestSuite) TestGetBatchByCommitmentID() {
	batchHash := utils.RandomHash()

	batch := &models.Batch{
		Hash:              batchHash,
		Type:              txtype.Transfer,
		ID:                models.MakeUint256(1),
		FinalisationBlock: 1234,
	}

	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitment := &models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeSignature(1, 2),
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

func (s *StateUpdateTestSuite) TestGetBatchByCommitmentID_NotExistentBatch() {
	commitment := &models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeSignature(1, 2),
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
			ID:                models.MakeUint256(1234),
			FinalisationBlock: 1234,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			ID:                models.MakeUint256(2000),
			FinalisationBlock: 2000,
		},
	}
	err := s.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.storage.AddBatch(&batches[1])
	s.NoError(err)

	actual, err := s.storage.GetLatestBatch()
	s.NoError(err)

	s.Equal(batches[1], *actual)
}

func (s *StateUpdateTestSuite) TestGetLatestBatch_NoBatches() {
	res, err := s.storage.GetLatestBatch()
	s.Equal(NewNotFoundError("batch"), err)
	s.Nil(res)
}



func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
