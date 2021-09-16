package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

type DisputeCT2TransitionTestSuite struct {
	DisputeTransitionTestSuite
}

func (s *DisputeCT2TransitionTestSuite) SetupTest() {
	s.DisputeTransitionTestSuite.SetupTest(txtype.Create2Transfer)
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_RemovesInvalidBatch() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Create2Transfer{
		{
			testutils.MakeCreate2Transfer(0, ref.Uint32(3), 0, 100, wallets[2].PublicKey()),
			testutils.MakeCreate2Transfer(1, ref.Uint32(4), 0, 100, wallets[0].PublicKey()),
		},
		{
			testutils.MakeCreate2Transfer(2, ref.Uint32(5), 0, 50, wallets[0].PublicKey()),
			testutils.MakeCreate2Transfer(2, ref.Uint32(6), 1, 500, wallets[0].PublicKey()),
		},
	}

	pubKeyIDs := [][]uint32{{3, 4}, {5, 6}}
	proofs := s.getStateMerkleProofs(commitmentTxs, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidBatch(commitmentTxs, pubKeyIDs, commitmentTxs[1][1].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 1)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[0], 1, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[0].ID)
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_FirstCommitment() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	commitmentTxs := [][]models.Create2Transfer{
		{
			testutils.MakeCreate2Transfer(0, ref.Uint32(4), 0, 500, wallets[1].PublicKey()),
		},
	}
	pubKeyIDs := [][]uint32{{4}}

	transfer := testutils.MakeCreate2Transfer(0, nil, 0, 50, wallets[1].PublicKey())
	submitC2TBatch(s.Assertions, s.client, s.rollupCtx, &transfer)

	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()

	pubKeyID, err := s.client.RegisterAccount(wallets[1].PublicKey(), registrations)
	s.NoError(err)
	s.EqualValues(4, *pubKeyID)

	proofs := s.getStateMerkleProofs(commitmentTxs, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()
	s.submitInvalidBatch(commitmentTxs, pubKeyIDs, commitmentTxs[0][0].Hash)

	remoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)
	s.Len(remoteBatches, 2)

	err = s.executionCtx.storage.MarkBatchAsSubmitted(&remoteBatches[0].Batch)
	s.NoError(err)

	err = s.disputeCtx.DisputeTransition(&remoteBatches[1], 0, proofs)
	s.NoError(err)

	s.checkBatchAfterDispute(remoteBatches[1].ID)
}

func (s *DisputeCT2TransitionTestSuite) TestDisputeTransition_ValidBatch() {
	wallets := setUserStates(s.Assertions, s.executionCtx, testDomain)

	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(0, nil, 0, 50, wallets[1].PublicKey()),
		testutils.MakeCreate2Transfer(0, ref.Uint32(4), 1, 100, wallets[1].PublicKey()),
	}
	pubKeyIDs := [][]uint32{{4}}

	submitC2TBatch(s.Assertions, s.client, s.rollupCtx, &transfers[0])

	proofs := s.getStateMerkleProofs([][]models.Create2Transfer{{transfers[1]}}, pubKeyIDs)

	s.beginTransaction()
	defer s.commitTransaction()

	transfers[1].ToStateID = nil
	submitC2TBatch(s.Assertions, s.client, s.rollupCtx, &transfers[1])

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

func (s *DisputeCT2TransitionTestSuite) getStateMerkleProofs(
	txs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
) []models.StateMerkleProof {
	feeReceiverStateID := uint32(0)

	s.beginTransaction()
	defer s.executionCtx.Rollback(nil)

	var stateProofs []models.StateMerkleProof
	var err error
	for i := range txs {
		_, stateProofs, err = s.executionCtx.ApplyCreate2TransfersForSync(txs[i], pubKeyIDs[i], feeReceiverStateID)
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

func (s *DisputeCT2TransitionTestSuite) submitInvalidBatch(
	txs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
	invalidTxHash common.Hash,
) *models.Batch {
	for i := range txs {
		stateIDs := s.resetToStateID(txs[i])
		err := s.executionCtx.storage.BatchAddCreate2Transfer(txs[i])
		s.NoError(err)
		s.setToStateID(txs[i], stateIDs)
	}

	pendingBatch, err := s.executionCtx.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)

	commitments := s.createInvalidCommitments(txs, pubKeyIDs, invalidTxHash)
	s.Len(commitments, len(txs))

	err = s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *DisputeCT2TransitionTestSuite) resetToStateID(txs []models.Create2Transfer) []*uint32 {
	stateIDs := make([]*uint32, 0, len(txs))
	for i := range txs {
		stateIDs = append(stateIDs, txs[i].ToStateID)
		txs[i].ToStateID = nil
	}
	return stateIDs
}

func (s *DisputeCT2TransitionTestSuite) setToStateID(txs []models.Create2Transfer, toStateIDs []*uint32) {
	for i := range txs {
		txs[i].ToStateID = toStateIDs[i]
	}
}

func (s *DisputeCT2TransitionTestSuite) createInvalidCommitments(
	commitmentTxs [][]models.Create2Transfer,
	pubKeyIDs [][]uint32,
	invalidTxHash common.Hash,
) []models.Commitment {
	commitmentID, err := s.executionCtx.createCommitmentID()
	s.NoError(err)

	commitments := make([]models.Commitment, 0, len(commitmentTxs))
	for i := range commitmentTxs {
		commitmentID.IndexInBatch = uint8(i)
		txs := commitmentTxs[i]
		combinedFee := models.MakeUint256(0)
		for j := range txs {
			receiverLeaf := newUserLeaf(*txs[j].ToStateID, pubKeyIDs[i][j], models.MakeUint256(0))
			combinedFee = s.applyTransfer(&txs[j], invalidTxHash, combinedFee, receiverLeaf)
		}
		if combinedFee.CmpN(0) > 0 {
			_, err := s.executionCtx.ApplyFee(0, combinedFee)
			s.NoError(err)
		}

		applyResult := &ApplyC2TForCommitmentResult{
			appliedTxs:     txs,
			addedPubKeyIDs: pubKeyIDs[i],
		}
		commitment, err := s.rollupCtx.buildTxCommitment(applyResult, commitmentID, 0, testDomain)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}

	return commitments
}

func newUserLeaf(stateID, pubKeyID uint32, tokenID models.Uint256) *models.StateLeaf {
	return &models.StateLeaf{
		StateID: stateID,
		UserState: models.UserState{
			PubKeyID: pubKeyID,
			TokenID:  tokenID,
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	}
}

func TestDisputeCT2TransitionTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeCT2TransitionTestSuite))
}
