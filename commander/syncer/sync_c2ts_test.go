package syncer

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/suite"
)

type SyncCreate2TransfersTestSuite struct {
	syncTxsTestSuite
}

func (s *SyncCreate2TransfersTestSuite) SetupTest() {
	s.syncTxsTestSuite.SetupTest(txtype.Create2Transfer)
}

func (s *SyncCreate2TransfersTestSuite) TestSyncTxs_AllValid() {
	input := s.generateValidTxs(3, 4)

	syncedTxs, stateProofs, err := s.syncCtx.SyncTxs(input, s.feeReceiverStateID)
	s.NoError(err)
	s.Len(syncedTxs, 3)
	s.Len(stateProofs, 7)
}

func (s *SyncCreate2TransfersTestSuite) TestSyncTxs_InvalidTransfer() {
	validC2Ts := s.generateValidTxs(2, 4)
	invalidC2Ts := s.generateInvalidTxs(3, 6)
	input := &SyncedC2Ts{
		txs:       append(validC2Ts.txs, invalidC2Ts.txs...),
		pubKeyIDs: append(validC2Ts.pubKeyIDs, invalidC2Ts.pubKeyIDs...),
	}

	syncedTxs, _, err := s.syncCtx.SyncTxs(input, s.feeReceiverStateID)
	s.Nil(syncedTxs)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *SyncCreate2TransfersTestSuite) TestSyncTxs_AppliesFee() {
	input := s.generateValidTxs(3, 4)

	_, _, err := s.syncCtx.SyncTxs(input, s.feeReceiverStateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiverStateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1030), feeReceiverState.Balance)
}

func (s *SyncCreate2TransfersTestSuite) TestSyncTxs_ReturnsCorrectStateProofsForZeroFee() {
	input := s.generateValidTxs(2, 5)
	for i := range input.txs {
		input.txs[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.SyncTxs(input, s.feeReceiverStateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *SyncCreate2TransfersTestSuite) TestSyncTxs_InvalidFeeReceiverTokenID() {
	feeReceiverStateID := uint32(4)
	_, err := s.storage.StateTree.Set(feeReceiverStateID, &models.UserState{
		PubKeyID: feeReceiverStateID,
		TokenID:  models.MakeUint256(4),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	input := s.generateValidTxs(2, 5)

	syncedTxs, _, err := s.syncCtx.SyncTxs(input, feeReceiverStateID)
	s.Nil(syncedTxs)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func (s *SyncCreate2TransfersTestSuite) generateValidTxs(txsAmount, startPubKeyID uint32) *SyncedC2Ts {
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

func (s *SyncCreate2TransfersTestSuite) generateInvalidTxs(txsAmount, startPubKeyID uint32) *SyncedC2Ts {
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

func TestSyncCreate2TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SyncCreate2TransfersTestSuite))
}
