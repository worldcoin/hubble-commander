package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.Storage
	db         *postgres.TestDB
	tree       *st.StateTree
	commitment models.Commitment
	batch      models.Batch
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchesTestSuite) SetupTest() {
	testDB, err := postgres.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage, nil}
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)

	hash := utils.RandomHash()
	s.commitment = commitment
	s.commitment.IncludedInBatch = &hash
	s.commitment.AccountTreeRoot = &hash

	s.batch = models.Batch{
		Hash:              hash,
		Type:              txtype.Transfer,
		FinalisationBlock: 42000,
	}
}

func (s *GetBatchesTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	err := addLeaf(s.storage, s.tree)
	s.NoError(err)
	err = s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 1)
	s.Equal(s.batch, result[0].Batch)
	s.Equal(getSubmissionBlock(s.batch.FinalisationBlock), result[0].SubmissionBlock)
}

func (s *GetBatchesTestSuite) TestGetBatchesByHash() {
	err := addLeaf(s.storage, s.tree)
	s.NoError(err)
	err = s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 1)
	s.Equal(s.batch, result[0].Batch)
	s.Equal(getSubmissionBlock(s.batch.FinalisationBlock), result[0].SubmissionBlock)
}

func (s *GetBatchesTestSuite) TestGetBatchesByHash_NoBatches() {
	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 0)
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
