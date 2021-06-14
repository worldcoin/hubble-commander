package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	senderState = models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(400),
		Nonce:      models.MakeUint256(0),
	}
	receiverState = models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}
	transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}
)

type ApplyTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *storage.TestStorage
	tree                *storage.StateTree
	transactionExecutor *transactionExecutor
}

func (s *ApplyTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransferTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = storage.NewStateTree(s.storage.Storage)
	s.transactionExecutor = newTestTransactionExecutor(s.storage.Storage, nil, config.GetTestConfig().Rollup, transactionExecutorOpts{})

	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
	}
	for i := range accounts {
		err = s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}
}

func (s *ApplyTransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyTransferTestSuite) TestCalculateStateAfterTransfer_UpdatesStates() {
	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(
		&senderState,
		&receiverState,
		&transfer,
	)
	s.NoError(err)

	s.Equal(models.MakeUint256(1), newSenderState.Nonce)
	s.Equal(models.MakeUint256(290), newSenderState.Balance)

	s.Equal(models.MakeUint256(0), newReceiverState.Nonce)
	s.Equal(models.MakeUint256(100), newReceiverState.Balance)

	s.NotEqual(&newSenderState, &senderState)
	s.NotEqual(&newReceiverState, &receiverState)
}

func (s *ApplyTransferTestSuite) TestCalculateStateAfterTransfer_Validation_Balance() {
	transferAboveBalance := transfer
	transferAboveBalance.Amount = models.MakeUint256(410)

	_, _, err := CalculateStateAfterTransfer(&senderState, &receiverState, &transferAboveBalance)
	s.Equal(ErrBalanceTooLow, err)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_Validation_ToStateID() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&create2Transfer, models.MakeUint256(1))
	s.NoError(transferError)
	s.Equal(ErrNilReceiverStateID, appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_Validation_Nonce() {
	transferWithBadNonce := transfer
	transferWithBadNonce.Nonce = models.MakeUint256(1)
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transferWithBadNonce, models.MakeUint256(1))
	s.Equal(ErrNonceTooHigh, transferError)
	s.NoError(appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_Validation_TokenIndex() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transfer, models.MakeUint256(3))
	s.Equal(appError, ErrIncorrectTokenIndices)
	s.NoError(transferError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.GetStateLeaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) setUserStatesInTree() {
	senderStateID := senderState.PubKeyID
	receiverStateID := receiverState.PubKeyID

	err := s.tree.Set(senderStateID, &senderState)
	s.NoError(err)
	err = s.tree.Set(receiverStateID, &receiverState)
	s.NoError(err)
}

func (s *ApplyTransferTestSuite) TestApplyFee() {
	receiverStateID := receiverState.PubKeyID
	err := s.tree.Set(receiverStateID, &receiverState)
	s.NoError(err)

	s.transactionExecutor.cfg.FeeReceiverPubKeyID = receiverStateID

	feeReceiverStateID, err := s.transactionExecutor.ApplyFee(models.MakeUint256(1), models.MakeUint256(555))
	s.NoError(err)
	s.Equal(receiverStateID, *feeReceiverStateID)

	receiverLeaf, err := s.storage.GetStateLeaf(receiverStateID)
	s.NoError(err)

	s.Equal(uint64(555), receiverLeaf.Balance.Uint64())
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
