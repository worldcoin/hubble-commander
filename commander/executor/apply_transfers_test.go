package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyTransfersTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *storage.TestStorage
	cfg         *config.RollupConfig
	rollupCtx   *RollupContext
	syncCtx     *SyncContext
	feeReceiver *FeeReceiver
}

func (s *ApplyTransfersTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransfersTestSuite) SetupTest() {
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
	s.rollupCtx = NewTestRollupContext(executionCtx, batchtype.Transfer)
	s.syncCtx = NewTestSyncContext(executionCtx, txtype.Transfer)

	s.feeReceiver = &FeeReceiver{
		StateID: 3,
		TokenID: models.MakeUint256(1),
	}
}

func (s *ApplyTransfersTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_AllValid() {
	generatedTransfers := generateValidTransfers(3)

	applyTxsResult, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(applyTxsResult.AppliedTxs(), 3)
	s.Len(applyTxsResult.InvalidTxs(), 0)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_SomeValid() {
	generatedTransfers := generateValidTransfers(2)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(3)...)

	applyTxsResult, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(applyTxsResult.AppliedTxs(), 2)
	s.Len(applyTxsResult.InvalidTxs(), 3)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_AppliesNoMoreThanLimit() {
	generatedTransfers := generateValidTransfers(13)

	applyTxsResult, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(applyTxsResult.AppliedTxs(), 6)
	s.Len(applyTxsResult.InvalidTxs(), 0)

	state, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(6), state.Nonce)
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_SavesTransferErrors() {
	generatedTransfers := generateValidTransfers(3)
	generatedTransfers = append(generatedTransfers, generateInvalidTransfers(2)...)

	for i := range generatedTransfers {
		err := s.storage.AddTransfer(&generatedTransfers[i])
		s.NoError(err)
	}

	applyTxsResult, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	s.Len(applyTxsResult.AppliedTxs(), 3)
	s.Len(applyTxsResult.InvalidTxs(), 2)

	for i := range generatedTransfers {
		transfer, err := s.storage.GetTransfer(generatedTransfers[i].Hash)
		s.NoError(err)
		if i < 3 {
			s.Nil(transfer.ErrorMessage)
		} else {
			s.Equal(*transfer.ErrorMessage, applier.ErrNonceTooLow.Error())
		}
	}
}

func (s *ApplyTransfersTestSuite) TestApplyTxs_AppliesFee() {
	generatedTransfers := generateValidTransfers(3)

	_, err := s.rollupCtx.ApplyTxs(generatedTransfers, s.cfg.MaxTxsPerCommitment, s.feeReceiver)
	s.NoError(err)

	feeReceiverState, err := s.rollupCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersForSync_AllValid() {
	transfers := generateValidTransfers(3)

	appliedTransfers, stateProofs, err := s.syncCtx.ApplyTransfersForSync(transfers, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(appliedTransfers, 3)
	s.Len(stateProofs, 7)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersForSync_InvalidTransfer() {
	transfers := generateValidTransfers(2)
	transfers = append(transfers, generateInvalidTransfers(2)...)

	appliedTransfers, _, err := s.syncCtx.ApplyTransfersForSync(transfers, s.feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Len(disputableErr.Proofs, 6)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersForSync_AppliesFee() {
	transfers := generateValidTransfers(3)

	_, _, err := s.syncCtx.ApplyTransfersForSync(transfers, s.feeReceiver.StateID)
	s.NoError(err)

	feeReceiverState, err := s.rollupCtx.storage.StateTree.Leaf(s.feeReceiver.StateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(1003), feeReceiverState.Balance)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersForSync_ReturnsCorrectStateProofsForZeroFee() {
	transfers := generateValidTransfers(2)
	for i := range transfers {
		transfers[i].Fee = models.MakeUint256(0)
	}

	_, stateProofs, err := s.syncCtx.ApplyTransfersForSync(transfers, s.feeReceiver.StateID)
	s.NoError(err)
	s.Len(stateProofs, 5)
}

func (s *ApplyTransfersTestSuite) TestApplyTransfersForSync_InvalidFeeReceiverTokenID() {
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

	transfers := generateValidTransfers(2)

	appliedTransfers, _, err := s.syncCtx.ApplyTransfersForSync(transfers, feeReceiver.StateID)
	s.Nil(appliedTransfers)

	var disputableErr *DisputableError
	s.ErrorAs(err, &disputableErr)
	s.Equal(Transition, disputableErr.Type)
	s.Equal(applier.ErrInvalidFeeReceiverTokenID.Error(), disputableErr.Reason)
	s.Len(disputableErr.Proofs, 5)
}

func generateValidTransfers(transfersAmount uint32) models.TransferArray {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := 0; i < int(transfersAmount); i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(uint64(i)),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func generateInvalidTransfers(transfersAmount uint64) []models.Transfer {
	transfers := make([]models.Transfer, 0, transfersAmount)
	for i := uint64(0); i < transfersAmount; i++ {
		transfer := models.Transfer{
			TransactionBase: models.TransactionBase{
				Hash:        utils.RandomHash(),
				TxType:      txtype.Transfer,
				FromStateID: 1,
				Amount:      models.MakeUint256(1_000_000),
				Fee:         models.MakeUint256(1),
				Nonce:       models.MakeUint256(0),
			},
			ToStateID: 2,
		}
		transfers = append(transfers, transfer)
	}
	return transfers
}

func TestApplyTransfersTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransfersTestSuite))
}
