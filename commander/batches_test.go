package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BatchTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	cfg     *config.RollupConfig
}

func (s *BatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *BatchTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	}
}

func (s *BatchTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *BatchTestSuite) TestSubmitBatch_ReturnsErrorWhenThereAreNotEnoughCommitments() {
	err := SubmitBatch(s.storage, nil, s.cfg)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *BatchTestSuite) TestSubmitBatch_ReturnsErrorWhenThereAreNotEnoughPendingCommitments() {
	batch := models.Batch{Hash: utils.RandomHash()}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	commitment := &models.Commitment{
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
		IncludedInBatch:   &batch.Hash,
	}
	_, err = s.storage.AddCommitment(commitment)
	s.NoError(err)

	err = SubmitBatch(s.storage, nil, s.cfg)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
