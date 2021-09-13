package executor

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

type SubmitTransferBatchTestSuite struct {
	SubmitBatchTestSuite
}

func (s *SubmitTransferBatchTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTest(txtype.Transfer)
	setupUser(&s.SubmitBatchTestSuite)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_ErrorsIfNotEnoughCommitments() {
	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
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

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitment := baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(txtype.Transfer, batch.Type)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_AddsCommitments() {
	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
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
