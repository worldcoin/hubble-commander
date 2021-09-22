package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/suite"
)

type SyncApplyCreate2TransfersTestSuite struct {
	SyncApplyTxsTestSuite
}

func (s *SyncApplyCreate2TransfersTestSuite) SetupTest() {
	s.SyncApplyTxsTestSuite.SetupTest(batchtype.Create2Transfer)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxs_AllValid() {
	input := s.generateValidTxs(3, 4)

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxs_InvalidTransfer() {
	validC2Ts := s.generateValidTxs(2, 4)
	invalidC2Ts := s.generateInvalidTxs(3, 6)
	input := &SyncedC2Ts{
		txs:       append(validC2Ts.txs, invalidC2Ts.txs...),
		pubKeyIDs: append(validC2Ts.pubKeyIDs, invalidC2Ts.pubKeyIDs...),
	}

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxs_AppliesFee() {
	input := s.generateValidTxs(3, 4)

	_, _, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1030), feeReceiverState.Balance)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxs_ReturnsCorrectStateProofsForZeroFee() {
	input := s.generateValidTxs(2, 5)
	for i := range input.txs {
		input.txs[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.ApplyTxs(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *SyncApplyCreate2TransfersTestSuite) TestApplyTxs_InvalidFeeReceiverTokenID() {
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

	input := s.generateValidTxs(2, 5)

	appliedTransfers, _, err := s.syncCtx.ApplyTxs(input, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func (s *SyncApplyCreate2TransfersTestSuite) generateValidTxs(txsAmount, startPubKeyID uint32) *SyncedC2Ts {
	syncedC2Ts := &SyncedC2Ts{
		txs:       make([]models.Create2Transfer, 0, txsAmount),
		pubKeyIDs: make([]uint32, 0, txsAmount),
	}

	for i := 0; i < int(txsAmount); i++ {
		tx := testutils.MakeCreate2Transfer(1, ref.Uint32(startPubKeyID), uint64(i), 1, &models.PublicKey{1, 2, 3})
		syncedC2Ts.txs = append(syncedC2Ts.txs, tx)
		syncedC2Ts.pubKeyIDs = append(syncedC2Ts.pubKeyIDs, startPubKeyID)
		startPubKeyID++
	}
	return syncedC2Ts
}

func (s *SyncApplyCreate2TransfersTestSuite) generateInvalidTxs(txsAmount, startPubKeyID uint32) *SyncedC2Ts {
	syncedC2Ts := &SyncedC2Ts{
		txs:       make([]models.Create2Transfer, 0, txsAmount),
		pubKeyIDs: make([]uint32, 0, txsAmount),
	}

	for i := 0; i < int(txsAmount); i++ {
		tx := testutils.MakeCreate2Transfer(1, ref.Uint32(startPubKeyID), uint64(i), 1_000_000, &models.PublicKey{1, 2, 3})
		syncedC2Ts.txs = append(syncedC2Ts.txs, tx)
		syncedC2Ts.pubKeyIDs = append(syncedC2Ts.pubKeyIDs, startPubKeyID)
		startPubKeyID++
	}
	return syncedC2Ts
}

func TestSyncApplyCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SyncApplyCreate2TransfersTestSuite))
}
