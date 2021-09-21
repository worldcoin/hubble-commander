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
