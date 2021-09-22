package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/stretchr/testify/suite"
)

type SyncApplyTransfersTestSuite struct {
	SyncApplyTxsTestSuite
}

func (s *SyncApplyTransfersTestSuite) SetupTest() {
	s.SyncApplyTxsTestSuite.SetupTest(batchtype.Transfer)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxs_AllValid() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(3),
	}

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxs_InvalidTransfer() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(2),
	}
	input.txs = append(input.txs, generateInvalidTransfers(2)...)

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxs_AppliesFee() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(3),
	}

	_, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxs_ReturnsCorrectStateProofsForZeroFee() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(2),
	}
	for i := range input.txs {
		input.txs[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxs_InvalidFeeReceiverTokenID() {
	feeReceiver := &FeeReceiver{
		StateID: 4,
		TokenID: models.MakeUint256(4),
	}
	_, err := s.storage.StateTree.Set(feeReceiver.StateID, &models.UserState{
		PubKeyID: 4,
		TokenID:  feeReceiver.TokenID,
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	input := &SyncedTransfers{
		txs: generateValidTransfers(2),
	}

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func TestSyncApplyTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SyncApplyTransfersTestSuite))
}
