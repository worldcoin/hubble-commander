package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputeTransferTransitionTestSuite struct {
	DisputeTransitionTestSuite
}

func (s *DisputeTransferTransitionTestSuite) SetupTest() {
	s.DisputeTransitionTestSuite.SetupTest(txtype.Transfer)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

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

	err = s.disputeCtx.DisputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Transfer{
		{
			testutils.MakeTransfer(0, 2, 0, 500),
		},
	}

	transfer := testutils.MakeTransfer(0, 2, 0, 50)
	submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfer)

	proofs := s.getStateMerkleProofs(commitmentTxs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidBatch(commitmentTxs, commitmentTxs[0][0].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *DisputeTransferTransitionTestSuite) TestDisputeTransition_ValidBatch() {
	setUserStates(s.Assertions, s.executionCtx, testDomain)

	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 2, 0, 50),
		testutils.MakeTransfer(0, 2, 1, 100),
	}

	submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfers[0])

	proofs := s.getStateMerkleProofs([][]models.Transfer{{transfers[1]}})

	s.beginTransaction()
	defer s.commitTransaction()
	submitTransferBatch(s.Assertions, s.client, s.rollupCtx, &transfers[1])

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)
	_, err = s.client.GetBatch(&remoteBatches[1].ID)
	s.NoError(err)
}

func (s *DisputeTransferTransitionTestSuite) getStateMerkleProofs(txs [][]models.Transfer) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.executionCtx.Rollback(nil)

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		_, stateProofs, err = s.executionCtx.ApplyTransfersForSync(txs[i], feeReceiverStateID)
		if err != nil {
			var disputableErr *DisputableError
			s.ErrorAs(err, &disputableErr)
			s.Equal(Transition, disputableErr.Type)
			s.Len(disputableErr.Proofs, len(txs[i])*2)
			return disputableErr.Proofs
		}
	}

	return stateProofs
}

func (s *DisputeTransferTransitionTestSuite) submitInvalidBatch(txs [][]models.Transfer, invalidTxHash common.Hash) *models.Batch {
	for i := range txs {
		err := s.executionCtx.storage.BatchAddTransfer(txs[i])
		s.NoError(err)
	}

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Transfer)
	s.NoError(err)

	commitments := s.createInvalidCommitments(txs, invalidTxHash)
	s.Len(commitments, len(txs))

	err = s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *DisputeTransferTransitionTestSuite) createInvalidCommitments(
	commitmentTxs [][]models.Transfer,
	invalidTxHash common.Hash,
) []models.Commitment {
	commitmentID, err := s.executionCtx.nextCommitmentID()
	s.NoError(err)

	commitments := make([]models.Commitment, 0, len(commitmentTxs))
	for i := range commitmentTxs {
		commitmentID.IndexInBatch = uint8(i)
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			receiverLeaf, err := s.executionCtx.storage.StateTree.Leaf(txs[j].ToStateID)
			s.NoError(err)
			combinedFee = s.applyTransfer(&txs[j], invalidTxHash, combinedFee, receiverLeaf)
		}
		if combinedFee.CmpN(0) > 0 {
			_, err := s.executionCtx.ApplyFee(0, combinedFee)
			s.NoError(err)
		}

		applyResult := &ApplyTransfersForCommitmentResult{
			appliedTxs: txs,
		}
		commitment, err := s.rollupCtx.buildCommitment(applyResult, commitmentID, 0, testDomain)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}

	return commitments
}

func TestDisputeTransferTransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeTransferTransitionTestSuite))
}
