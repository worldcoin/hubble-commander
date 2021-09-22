package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SyncApplyTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *storage.TestStorage
	cfg         *config.RollupConfig
	syncCtx     *SyncContext
	feeReceiver *FeeReceiver
}

func (s *SyncApplyTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SyncApplyTransfersTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		FeeReceiverPubKeyID: 3,
		MaxTxsPerCommitment: 6,
	}

	senderState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	receiverState := models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	feeReceiverState := models.UserState{
		PubKeyID: 3,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, &senderState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(2, &receiverState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(3, &feeReceiverState)
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, nil, s.cfg)
	s.syncCtx = NewTestSyncContext(executionCtx, batchtype.Transfer)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *SyncApplyTransfersTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxsForSync_AllValid() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(3),
	}

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxsForSync_InvalidTransfer() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(2),
	}
	input.txs = append(input.txs, generateInvalidTransfers(2)...)

	appliedTransfers, _, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxsForSync_AppliesFee() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(3),
	}

	_, _, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.syncCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxsForSync_ReturnsCorrectStateProofsForZeroFee() {
	input := &SyncedTransfers{
		txs: generateValidTransfers(2),
	}
	for i := range input.txs {
		input.txs[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.ApplyTxsForSync(input, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *SyncApplyTransfersTestSuite) TestApplyTxsForSync_InvalidFeeReceiverTokenID() {
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

	appliedTransfers, _, err := s.syncCtx.ApplyTxsForSync(input, feeReceiver.StateID)
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
