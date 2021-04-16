package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
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

func (s *BatchTestSuite) Test_AddBatch_AddAndRetrieve() {
	batch := &models.Batch{
		Hash:              utils.RandomHash(),
		ID:                models.MakeUint256(1),
		FinalisationBlock: 1234,
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatch(batch.Hash)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *BatchTestSuite) Test_GetBatchByID() {
	batch := &models.Batch{
		Hash:              utils.RandomHash(),
		ID:                models.MakeUint256(1234),
		FinalisationBlock: 1234,
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	actual, err := s.storage.GetBatchByID(batch.ID)
	s.NoError(err)

	s.Equal(batch, actual)
}

func (s *StateUpdateTestSuite) Test_GetBatch_NonExistentBatch() {
	res, err := s.storage.GetBatch(common.Hash{1, 2, 3, 4})
	s.Equal(NewNotFoundErr("batch"), err)
	s.Nil(res)
}

func (s *StateUpdateTestSuite) Test_GetBatchById_NonExistentBatch() {
	res, err := s.storage.GetBatchByID(models.MakeUint256(42))
	s.Equal(NewNotFoundErr("batch"), err)
	s.Nil(res)
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
