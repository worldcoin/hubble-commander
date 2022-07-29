package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type SubmitMassMigrationBatchTestSuite struct {
	submitBatchTestSuite
	commitment *models.MMCommitmentWithTxs
}

func (s *SubmitMassMigrationBatchTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTest(batchtype.MassMigration)
	s.setupUser()

	s.commitment = &models.MMCommitmentWithTxs{
		MMCommitment: models.MMCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchtype.MassMigration,
				PostStateRoot: utils.RandomHash(),
			},
			CombinedSignature: models.MakeRandomSignature(),
			BodyHash:          utils.NewRandomHash(),
			Meta: &models.MassMigrationMeta{
				SpokeID:     1,
				TokenID:     models.MakeUint256(1),
				Amount:      models.MakeUint256(10),
				FeeReceiver: 1,
			},
			WithdrawRoot: utils.RandomHash(),
		},
		Transactions: utils.RandomBytes(8),
	}
	s.commitment.Type = batchtype.MassMigration
}

func (s *SubmitMassMigrationBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.MassMigration)
	s.NoError(err)

	s.commitment.ID.BatchID = pendingBatch.ID

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, []models.CommitmentWithTxs{s.commitment})
	s.NoError(err)

	s.client.GetBackend().Commit()

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitMassMigrationBatchTestSuite) TestSubmitBatch_StoresPendingBatchRecord() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.MassMigration)
	s.NoError(err)

	s.commitment.ID.BatchID = pendingBatch.ID

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, []models.CommitmentWithTxs{s.commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(batchtype.MassMigration, batch.Type)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitMassMigrationBatchTestSuite) TestSubmitBatch_AddsCommitments() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.MassMigration)
	s.NoError(err)

	commitments := getMMCommitments(2, pendingBatch.ID)

	err = s.txsCtx.SubmitBatch(context.Background(), pendingBatch, commitments)
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for i := range commitments {
		commit, err := s.storage.GetCommitment(&commitments[i].ToMMCommitmentWithTxs().ID)
		s.NoError(err)
		s.Equal(commitments[i].ToMMCommitmentWithTxs().MMCommitment, *commit.ToMMCommitment())
		s.Equal(batch.ID, commit.ToMMCommitment().ID.BatchID)
	}
}

func TestSubmitMassMigrationBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitMassMigrationBatchTestSuite))
}
