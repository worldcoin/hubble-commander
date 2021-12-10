package disputer

import (
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputeMMTransitionTestSuite struct {
	disputeTransitionTestSuite
}

func (s *DisputeMMTransitionTestSuite) SetupTest() {
	s.SetupTestWithConfig(batchtype.MassMigration, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    2,
		DisableSignatures:      true,
	})
}

func (s *DisputeMMTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	txs := [][]models.MassMigration{
		{
			testutils.MakeMassMigration(0, 1, 0, 100),
			testutils.MakeMassMigration(1, 1, 0, 100),
		},
		{
			testutils.MakeMassMigration(2, 1, 0, 20),
			testutils.MakeMassMigration(2, 1, 1, 20),
		},
	}

	s.submitInvalidBatch(txs)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	proofs := s.getInvalidBatchStateProofs(remoteBatches[0])
	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 1, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeMMTransitionTestSuite) TestDisputeTransition_ValidBatch() {
	tx := testutils.MakeMassMigration(0, 1, 0, 50)
	proofs := s.getValidBatchStateProofs([][]models.MassMigration{{tx}})

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitBatch(&tx)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[0].GetBase().ID)
	s.NoError(err)
}

func (s *DisputeMMTransitionTestSuite) getInvalidBatchStateProofs(remoteBatch eth.DecodedBatch) []models.StateMerkleProof {
	s.beginTransaction()
	defer s.rollback()

	err := s.syncCtx.SyncCommitments(remoteBatch)
	s.Error(err)

	var disputableErr *syncer.DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(syncer.Transition, disputableErr.Type)
	return disputableErr.Proofs
}

func (s *DisputeMMTransitionTestSuite) getValidBatchStateProofs(txs [][]models.MassMigration) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.rollback()

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		input := syncer.NewSyncedMMs(txs[i])
		_, stateProofs, err = s.syncCtx.SyncTxs(input, feeReceiverStateID)
		s.NoError(err)
	}

	return stateProofs
}

func (s *DisputeMMTransitionTestSuite) submitInvalidBatch(txs [][]models.MassMigration) *models.Batch {
	s.beginTransaction()
	defer s.rollback()
	for i := range txs {
		err := s.disputeCtx.storage.BatchAddMassMigration(txs[i])
		s.NoError(err)
	}

	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)
	fmt.Println(*pendingBatch.PrevStateRoot)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)

	batchData.Commitments()[batchData.Len()-1].PostStateRoot = common.Hash{1, 2, 3}

	err = s.txsCtx.SubmitBatch(pendingBatch, batchData)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func TestDisputeMMTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeMMTransitionTestSuite))
}
