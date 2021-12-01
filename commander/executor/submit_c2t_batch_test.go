package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SubmitC2TBatchTestSuite struct {
	submitBatchTestSuite
}

func (s *SubmitC2TBatchTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTest(batchtype.Create2Transfer)
	s.setupUser()
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(1), nextBatchID)

	commitment := baseCommitment
	commitment.ID.BatchID = models.MakeUint256FromBig(*nextBatchID)

	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)
	err = s.txsCtx.SubmitBatch(pendingBatch, &CreateTxCommitmentsResult{commitments: []models.CommitmentWithTxs{commitment}})
	s.NoError(err)

	s.client.GetBackend().Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitC2TBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Create2Transfer)
	s.NoError(err)

	commitment := baseCommitment
	commitment.ID.BatchID = pendingBatch.ID

	err = s.txsCtx.SubmitBatch(pendingBatch, &CreateTxCommitmentsResult{commitments: []models.CommitmentWithTxs{commitment}})
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
	commitments := getCommitments(2, pendingBatch.ID)

	err = s.txsCtx.SubmitBatch(pendingBatch, &CreateTxCommitmentsResult{commitments: commitments})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for i := range commitments {
		commit, err := s.storage.GetTxCommitment(&commitments[i].ID)
		s.NoError(err)
		s.Equal(commitments[i].TxCommitment, *commit)
		s.Equal(batch.ID, commit.ID.BatchID)
	}
}

func TestSubmitC2TBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitC2TBatchTestSuite))
}
