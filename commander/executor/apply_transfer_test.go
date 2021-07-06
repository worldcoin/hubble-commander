package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	senderState = models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(400),
		Nonce:    models.MakeUint256(0),
	}
	receiverState = models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
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
	transactionExecutor *TransactionExecutor
}

func (s *ApplyTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransferTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = storage.NewStateTree(s.storage.Storage)
	s.transactionExecutor = NewTestTransactionExecutor(
		s.storage.Storage,
		nil,
		&config.RollupConfig{FeeReceiverPubKeyID: 0},
		context.Background(),
	)

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
	newSenderState, newReceiverState, err := calculateStateAfterTransfer(
		senderState,
		receiverState,
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

func (s *ApplyTransferTestSuite) TestCalculateStateAfterTransfer_ValidatesTokenAmount() {
	invalidTransfer := transfer
	invalidTransfer.Amount = models.MakeUint256(0)
	_, _, err := calculateStateAfterTransfer(
		senderState,
		receiverState,
		&invalidTransfer,
	)
	s.Equal(ErrInvalidTokenAmount, err)
}

func (s *ApplyTransferTestSuite) TestCalculateStateAfterTransfer_ValidatesBalance() {
	transferAboveBalance := transfer
	transferAboveBalance.Amount = models.MakeUint256(410)

	_, _, err := calculateStateAfterTransfer(senderState, receiverState, &transferAboveBalance)
	s.Equal(ErrBalanceTooLow, err)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesToStateID() {
	c2T := create2Transfer
	transferError, appError := s.transactionExecutor.ApplyTransfer(&c2T, models.MakeUint256(1))
	s.NoError(transferError)
	s.Equal(ErrNilReceiverStateID, appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesSenderTokenID() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transfer, models.MakeUint256(3))
	s.NoError(transferError)
	s.Equal(appError, ErrInvalidSenderTokenID)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesReceiverTokenID() {
	s.setUserStatesInTree()

	receiverWithChangedToken := receiverState
	receiverWithChangedToken.TokenID = models.MakeUint256(2)
	_, err := s.tree.Set(2, &receiverWithChangedToken)
	s.NoError(err)

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transfer, models.MakeUint256(1))
	s.NoError(transferError)
	s.Equal(appError, ErrInvalidReceiverTokenID)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesNonce() {
	transferWithBadNonce := transfer
	transferWithBadNonce.Nonce = models.MakeUint256(1)
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transferWithBadNonce, models.MakeUint256(1))
	s.Equal(ErrNonceTooHigh, transferError)
	s.NoError(appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_UpdatesStatesCorrectly() {
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

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesToStateID() {
	c2T := create2Transfer
	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&c2T, models.MakeUint256(1))
	s.Nil(synced)
	s.NoError(transferError)
	s.Equal(ErrNilReceiverStateID, appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsSenderProofForCalculateStateAfterTransferValidations() {
	s.setUserStatesInTree()

	bigTransfer := transfer
	bigTransfer.Amount = models.MakeUint256(1_000_000)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&bigTransfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.Equal(ErrBalanceTooLow, transferError)
	s.NoError(appError)

	s.Equal(&bigTransfer, synced.transfer)
	s.Equal(senderState, synced.senderStateProof.UserState)
	s.Len(synced.senderStateProof.Witness, storage.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesSenderTokenID() {
	s.setUserStatesInTree()

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(3))
	s.NotNil(synced)
	s.Equal(ErrInvalidSenderTokenID, transferError)
	s.NoError(appError)

	s.Equal(&transfer, synced.transfer)
	s.Equal(senderState, synced.senderStateProof.UserState)
	s.Len(synced.senderStateProof.Witness, storage.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesReceiverTokenID() {
	s.setUserStatesInTree()

	receiverWithChangedToken := receiverState
	receiverWithChangedToken.TokenID = models.MakeUint256(2)
	_, err := s.tree.Set(2, &receiverWithChangedToken)
	s.NoError(err)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.Equal(ErrInvalidReceiverTokenID, transferError)
	s.NoError(appError)

	s.Equal(&transfer, synced.transfer)
	s.Equal(senderState, synced.senderStateProof.UserState)
	s.Len(synced.senderStateProof.Witness, storage.StateTreeDepth)
	s.Equal(receiverWithChangedToken, synced.receiverStateProof.UserState)
	s.Len(synced.receiverStateProof.Witness, storage.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsTransferWithUpdatedNonce() {
	s.setUserStatesInTree()
	transferWithModifiedNonce := transfer
	transferWithModifiedNonce.Nonce = models.MakeUint256(1234)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transferWithModifiedNonce, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1234), transferWithModifiedNonce.Nonce)
	s.Equal(models.MakeUint256(0), synced.transfer.GetNonce())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_UpdatesStatesCorrectly() {
	s.setUserStatesInTree()

	_, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.storage.GetStateLeaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.GetStateLeaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsProofs() {
	s.setUserStatesInTree()

	sync, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(senderState, sync.senderStateProof.UserState)
	s.Len(sync.senderStateProof.Witness, storage.StateTreeDepth)
	s.Equal(receiverState, sync.receiverStateProof.UserState)
	s.Len(sync.receiverStateProof.Witness, storage.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_SetsNonce() {
	s.setUserStatesInTree()

	_, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	sync, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1), sync.transfer.GetNonce())
}

func (s *ApplyTransferTestSuite) setUserStatesInTree() {
	senderStateID := senderState.PubKeyID
	receiverStateID := receiverState.PubKeyID

	_, err := s.tree.Set(senderStateID, &senderState)
	s.NoError(err)
	_, err = s.tree.Set(receiverStateID, &receiverState)
	s.NoError(err)
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
