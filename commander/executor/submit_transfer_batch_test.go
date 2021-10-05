package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SubmitTransferBatchTestSuite struct {
	submitBatchTestSuite
}

func (s *SubmitTransferBatchTestSuite) SetupTest() {
	s.testSuiteWithRollupContext.SetupTest(batchtype.Transfer)
	s.setupUser()
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_ErrorsIfNotEnoughCommitments() {
	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)
	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{})
	s.Equal(ErrNotEnoughCommitments, err)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(1), nextBatchID)

	commitment := baseCommitment
	commitment.ID.BatchID = models.MakeUint256FromBig(*nextBatchID)

	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)
	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)

	commitment := baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(batchtype.Transfer, batch.Type)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_AddsCommitments() {
	pendingBatch, err := s.rollupCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)
	commitments := getCommitments(2, pendingBatch.ID)

	err = s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for i := range commitments {
		commit, err := s.storage.GetCommitment(&commitments[i].ID)
		s.NoError(err)
		s.Equal(commitments[i], *commit)
		s.Equal(batch.ID, commit.ID.BatchID)
	}
}

func TestSubmitTransferBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTransferBatchTestSuite))
}