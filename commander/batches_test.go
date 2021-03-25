package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
	}
)

type BatchTestSuite struct {
	*require.Assertions
	suite.Suite
	db         *db.TestDB
	storage    *storage.Storage
	cfg        *config.RollupConfig
	testClient *eth.TestClient
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

	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
}

func (s *BatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *BatchTestSuite) TestSubmitBatch_ReturnsErrorWhenThereAreNotEnoughCommitments() {
	err := SubmitBatch(s.storage, s.testClient.Client, s.cfg)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *BatchTestSuite) TestSubmitBatch_ReturnsErrorWhenThereAreNotEnoughPendingCommitments() {
	batch := models.Batch{Hash: utils.RandomHash()}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = &batch.Hash
	_, err = s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	err = SubmitBatch(s.storage, s.testClient.Client, s.cfg)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *BatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = SubmitBatch(s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	nextBatchID, err := s.testClient.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *BatchTestSuite) TestSubmitBatch_StoresBatchRecord() {
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = SubmitBatch(s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)
	s.NotNil(batch)
}

func (s *BatchTestSuite) TestSubmitBatch_MarksCommitmentsAsIncluded() {
	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)
	id2, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = SubmitBatch(s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range []int32{*id, *id2} {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.Hash, *commit.IncludedInBatch)
	}
}

func TestBatchTestSuite(t *testing.T) {
	suite.Run(t, new(BatchTestSuite))
}
