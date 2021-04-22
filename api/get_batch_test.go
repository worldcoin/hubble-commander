package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.Storage
	db         *db.TestDB
	commitment models.Commitment
	batch      models.Batch
}

func (s *GetBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage, nil}
	s.db = testDB

	hash := utils.RandomHash()
	s.commitment = commitment
	s.commitment.IncludedInBatch = &hash

	s.batch = models.Batch{
		Hash: hash,
		Type: txtype.Transfer,
	}
}

func (s *GetBatchTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetBatchTestSuite) TestGetBatchByHash() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch.Hash, result.Hash)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NoCommitments() {
	batch := models.Batch{
		Hash: *s.commitment.IncludedInBatch,
		Type: txtype.Transfer,
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 0)
	s.Equal(batch.Hash, result.Hash)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonExistentBatch() {
	result, err := s.api.GetBatchByHash(*s.commitment.IncludedInBatch)
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch.Hash, result.Hash)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NoCommitments() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 0)
	s.Equal(s.batch.Hash, result.Hash)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NonExistentBatch() {
	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
