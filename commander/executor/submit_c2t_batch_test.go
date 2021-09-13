package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SubmitC2TBatchTestSuite struct {
	SubmitBatchTestSuite
}

func (s *SubmitC2TBatchTestSuite) SetupTest() {
	s.TestSuiteWithRollupContext.SetupTest(txtype.Create2Transfer)
	setupUser(&s.SubmitBatchTestSuite)
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(1), nextBatchID)

	commitment := baseCommitment
	commitment.ID.BatchID = models.MakeUint256FromBig(*nextBatchID)

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)
	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)

	commitment := baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.rollupCtx.SubmitBatch(pendingBatch, []models.Commitment{commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(pendingBatch.Type, batch.Type)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_AddsCommitments() {
	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Create2Transfer)
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
