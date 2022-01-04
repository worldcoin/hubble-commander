package executor

import (
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
	commitment models.TxCommitmentWithTxs
	batchData  BatchData
}

func (s *SubmitMassMigrationBatchTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTest(batchtype.MassMigration)
	s.setupUser()

	s.commitment = baseCommitment
	s.commitment.Type = batchtype.MassMigration
	s.batchData = &MassMigrationBatchData{
		commitments: make([]models.TxCommitmentWithTxs, 0, 1),
		metas: []models.MassMigrationMeta{
			{
				SpokeID:     1,
				TokenID:     models.MakeUint256(1),
				Amount:      models.MakeUint256(10),
				FeeReceiver: 1,
			},
		},
		withdrawRoots: []common.Hash{utils.RandomHash()},
	}
}

func (s *SubmitMassMigrationBatchTestSuite) TestSubmitBatch_SubmitsCommitmentsOnChain() {
	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.MassMigration)
	s.NoError(err)

	s.commitment.ID.BatchID = pendingBatch.ID
	s.batchData.AddCommitment(&s.commitment)

	err = s.txsCtx.SubmitBatch(pendingBatch, s.batchData)
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
	s.batchData.AddCommitment(&s.commitment)

	err = s.txsCtx.SubmitBatch(pendingBatch, s.batchData)
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

	commitments := getCommitments(2, pendingBatch.ID, batchtype.MassMigration)
	s.batchData.AddCommitment(&commitments[0])
	s.batchData.AddCommitment(&commitments[1])
	s.batchData.AddMeta(&s.batchData.Metas()[0])
	s.batchData.AddWithdrawRoot(utils.RandomHash())

	err = s.txsCtx.SubmitBatch(pendingBatch, s.batchData)
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

func TestSubmitMassMigrationBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitMassMigrationBatchTestSuite))
}
