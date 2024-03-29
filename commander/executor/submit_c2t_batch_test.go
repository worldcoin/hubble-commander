package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SubmitC2TBatchTestSuite struct {
	submitBatchTestSuite
	baseCommitment models.TxCommitmentWithTxs
}

func (s *SubmitC2TBatchTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTest(batchtype.Create2Transfer)
	s.setupUser()

	s.baseCommitment = baseCommitment
	s.baseCommitment.Type = batchtype.Create2Transfer
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)

	commitment := s.baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, []models.CommitmentWithTxs{&commitment})
	s.NoError(err)

	s.client.GetBackend().Commit()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)

	commitment := s.baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, []models.CommitmentWithTxs{&commitment})
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
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)
	commitments := getTxCommitments(2, pendingBatch.ID, batchtype.Create2Transfer)

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, commitments)
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for i := range commitments {
		commit, err := s.storage.GetCommitment(&commitments[i].ToTxCommitmentWithTxs().ID)
		s.NoError(err)
		s.Equal(commitments[i].ToTxCommitmentWithTxs().TxCommitment, *commit.ToTxCommitment())
		s.Equal(batch.ID, commit.ToTxCommitment().ID.BatchID)
	}
}

func TestSubmitC2TBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitC2TBatchTestSuite))
}
