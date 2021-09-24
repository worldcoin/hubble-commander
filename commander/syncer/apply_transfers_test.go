package syncer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/suite"
)

type ApplyTransfersTestSuite struct {
	applyTxsTestSuite
}

func (s *ApplyTransfersTestSuite) SetupTest() {
	s.applyTxsTestSuite.SetupTest(batchtype.Transfer)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_AllValid() {
	input := &SyncedTransfers{
		txs: testutils.GenerateValidTransfers(3),
	}

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiverStateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_InvalidTransfer() {
	input := &SyncedTransfers{
		txs: testutils.GenerateValidTransfers(2),
	}
	input.txs = append(input.txs, testutils.GenerateInvalidTransfers(2)...)

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiverStateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_AppliesFee() {
	input := &SyncedTransfers{
		txs: testutils.GenerateValidTransfers(3),
	}

	_, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiverStateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiverStateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_ReturnsCorrectStateProofsForZeroFee() {
	input := &SyncedTransfers{
		txs: testutils.GenerateValidTransfers(2),
	}
	for i := range input.txs {
		input.txs[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiverStateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_InvalidFeeReceiverTokenID() {
	feeReceiverStateID := uint32(4)
	_, err := s.storage.StateTree.Set(feeReceiverStateID, &models.UserState{
		PubKeyID: feeReceiverStateID,
		TokenID:  models.MakeUint256(4),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	input := &SyncedTransfers{
		txs: testutils.GenerateValidTransfers(2),
	}

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, feeReceiverStateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func TestApplyTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransfersTestSuite))
}
