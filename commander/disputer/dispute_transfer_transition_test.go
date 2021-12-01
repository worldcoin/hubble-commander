package disputer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputeTransferTransitionTestSuite struct {
	disputeTransitionTestSuite
}

func (s *DisputeTransferTransitionTestSuite) SetupTest() {
	s.disputeTransitionTestSuite.SetupTest(batchtype.Transfer)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	setUserStates(s.Assertions, s.disputeCtx, &bls.TestDomain)

	commitmentTxs := [][]models.Transfer{
		{
			testutils.MakeTransfer(0, 2, 0, 100),
			testutils.MakeTransfer(1, 0, 0, 100),
		},
		{
			testutils.MakeTransfer(2, 0, 0, 50),
			testutils.MakeTransfer(2, 0, 1, 500),
		},
	}

	proofs := s.getStateMerkleProofs(commitmentTxs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidBatch(commitmentTxs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(remoteBatches[0].ToDecodedTxBatch(), 1, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[0].GetBase().ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	setUserStates(s.Assertions, s.disputeCtx, &bls.TestDomain)

	commitmentTxs := [][]models.Transfer{
		{
			testutils.MakeTransfer(0, 2, 0, 500),
		},
	}

	transfer := testutils.MakeTransfer(0, 2, 0, 50)
	s.submitBatch(&transfer)

	proofs := s.getStateMerkleProofs(commitmentTxs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidBatch(commitmentTxs, commitmentTxs[0][0].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.syncCtx.UpdateExistingBatch(remoteBatches[0], common.Hash{1, 2, 3})
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(remoteBatches[1].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)

	checkRemoteBatchAfterDispute(s.Assertions, s.client, &remoteBatches[1].GetBase().ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_ValidBatch() {
	setUserStates(s.Assertions, s.disputeCtx, &bls.TestDomain)

	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 1, 100),
	}

	s.submitBatch(&transfers[0])

	proofs := s.getStateMerkleProofs([][]models.Transfer{{transfers[1]}})

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitBatch(&transfers[1])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.syncCtx.UpdateExistingBatch(remoteBatches[0], common.Hash{1, 2, 3})
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(remoteBatches[1].ToDecodedTxBatch(), 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[1].GetBase().ID)
	s.NoError(err)
}

func (s *DisputeTransferTransitionTestSuite) getStateMerkleProofs(txs [][]models.Transfer) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.rollback()

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		input := syncer.NewSyncedTransfers(txs[i])
		_, stateProofs, err = s.syncCtx.SyncTxs(input, feeReceiverStateID)
		if err != nil {
			var disputableErr *syncer.DisputableError
			s.ErrorAs(err, &disputableErr)
			s.Equal(syncer.Transition, disputableErr.Type)
			s.Len(disputableErr.Proofs, len(txs[i])*2)
			return disputableErr.Proofs
		}
	}

	return stateProofs
}

func (s *DisputeTransferTransitionTestSuite) submitInvalidBatch(txs [][]models.Transfer, invalidTxHash common.Hash) *models.Batch {
	for i := range txs {
		err := s.disputeCtx.storage.BatchAddTransfer(txs[i])
		s.NoError(err)
	}

	pendingBatch, err := s.txsCtx.NewPendingBatch(batchtype.Transfer)
	s.NoError(err)

	result := s.createInvalidCommitments(txs, invalidTxHash)
	s.Equal(result.Len(), len(txs))

	err = s.txsCtx.SubmitBatch(pendingBatch, result)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *DisputeTransferTransitionTestSuite) createInvalidCommitments(
	commitmentTxs [][]models.Transfer,
	invalidTxHash common.Hash,
) executor.CreateCommitmentsResult {
	commitmentID, err := s.txsCtx.NextCommitmentID()
	s.NoError(err)

	result := s.txsCtx.Executor.NewCreateCommitmentsResult(uint32(len(commitmentTxs)))
	for i := range commitmentTxs {
		commitmentID.IndexInBatch = uint8(i)
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			receiverLeaf, err := s.disputeCtx.storage.StateTree.Leaf(txs[j].ToStateID)
			s.NoError(err)
			combinedFee = s.applyTransfer(&txs[j], invalidTxHash, combinedFee, receiverLeaf)
		}
		if combinedFee.CmpN(0) > 0 {
			_, err := s.txsCtx.Applier.ApplyFee(0, combinedFee)
			s.NoError(err)
		}

		executeTxsResult := s.txsCtx.Executor.NewExecuteTxsResult(uint32(len(txs)))
		for j := range txs {
			executeTxsResult.AddApplied(applier.NewApplySingleTransferResult(&txs[j]))
		}
		executeTxsForCommitmentResult := s.txsCtx.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult)
		commitment, err := s.txsCtx.BuildCommitment(executeTxsForCommitmentResult, commitmentID, 0)
		s.NoError(err)
		result.AddCommitment(commitment)
	}

	return result
}

func TestDisputeTransferTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransferTransitionTestSuite))
}
