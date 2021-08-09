package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
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
)

type ApplyTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	transactionExecutor *TransactionExecutor
	transfer            models.Transfer
	receiverLeaf        models.StateLeaf
}

func (s *ApplyTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}
}

func (s *ApplyTransferTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.transactionExecutor = NewTestTransactionExecutor(
		s.storage.Storage,
		nil,
		&config.RollupConfig{FeeReceiverPubKeyID: 0},
		context.Background(),
	)

	s.receiverLeaf = models.StateLeaf{
		StateID:   receiverState.PubKeyID,
		UserState: receiverState,
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
		&s.transfer,
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
	invalidTransfer := s.transfer
	invalidTransfer.Amount = models.MakeUint256(0)
	_, _, err := calculateStateAfterTransfer(
		senderState,
		receiverState,
		&invalidTransfer,
	)
	s.Equal(ErrInvalidTokenAmount, err)
}

func (s *ApplyTransferTestSuite) TestCalculateStateAfterTransfer_ValidatesBalance() {
	transferAboveBalance := s.transfer
	transferAboveBalance.Amount = models.MakeUint256(410)

	_, _, err := calculateStateAfterTransfer(senderState, receiverState, &transferAboveBalance)
	s.Equal(ErrBalanceTooLow, err)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesSenderTokenID() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&s.transfer, &s.receiverLeaf, models.MakeUint256(3))
	s.NoError(transferError)
	s.Equal(appError, ErrInvalidSenderTokenID)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesReceiverTokenID() {
	s.setUserStatesInTree()

	receiverWithChangedToken := s.receiverLeaf
	receiverWithChangedToken.TokenID = models.MakeUint256(2)

	transferError, appError := s.transactionExecutor.ApplyTransfer(&s.transfer, &receiverWithChangedToken, models.MakeUint256(1))
	s.NoError(transferError)
	s.Equal(appError, ErrInvalidReceiverTokenID)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesNonce() {
	transferWithBadNonce := s.transfer
	transferWithBadNonce.Nonce = models.MakeUint256(1)
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&transferWithBadNonce, &s.receiverLeaf, models.MakeUint256(1))
	s.Equal(ErrNonceTooHigh, transferError)
	s.NoError(appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_UpdatesStatesCorrectly() {
	s.setUserStatesInTree()

	transferError, appError := s.transactionExecutor.ApplyTransfer(&s.transfer, &s.receiverLeaf, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsSenderProofForCalculateStateAfterTransferValidations() {
	s.setUserStatesInTree()

	bigTransfer := s.transfer
	bigTransfer.Amount = models.MakeUint256(1_000_000)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&bigTransfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.Equal(ErrBalanceTooLow, transferError)
	s.NoError(appError)

	s.Equal(&bigTransfer, synced.Transfer)
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesSenderTokenID() {
	s.setUserStatesInTree()

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(3))
	s.NotNil(synced)
	s.Equal(ErrInvalidSenderTokenID, transferError)
	s.NoError(appError)

	s.Equal(&s.transfer, synced.Transfer)
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesReceiverTokenID() {
	s.setUserStatesInTree()

	receiverWithChangedToken := receiverState
	receiverWithChangedToken.TokenID = models.MakeUint256(2)
	_, err := s.storage.StateTree.Set(2, &receiverWithChangedToken)
	s.NoError(err)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.Equal(ErrInvalidReceiverTokenID, transferError)
	s.NoError(appError)

	s.Equal(&s.transfer, synced.Transfer)
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
	s.Equal(receiverWithChangedToken, *synced.ReceiverStateProof.UserState)
	s.Len(synced.ReceiverStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsTransferWithUpdatedNonce() {
	s.setUserStatesInTree()
	transferWithModifiedNonce := s.transfer
	transferWithModifiedNonce.Nonce = models.MakeUint256(1234)

	synced, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transferWithModifiedNonce, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1234), transferWithModifiedNonce.Nonce)
	s.Equal(models.MakeUint256(0), synced.Transfer.GetNonce())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_UpdatesStatesCorrectly() {
	s.setUserStatesInTree()

	_, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsProofs() {
	s.setUserStatesInTree()

	sync, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(senderState, *sync.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
	s.Equal(receiverState, *sync.ReceiverStateProof.UserState)
	s.Len(sync.ReceiverStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesNotExistingSenderState() {
	s.setUserStatesInTree()

	senderLeaf, err := s.storage.StateTree.LeafOrEmpty(10)
	s.NoError(err)

	transfer := s.transfer
	transfer.FromStateID = senderLeaf.StateID

	sync, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.ErrorIs(transferError, ErrBalanceTooLow)
	s.Equal(senderLeaf.UserState, *sync.Proofs.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_SetsNonce() {
	s.setUserStatesInTree()

	_, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	sync, transferError, appError := s.transactionExecutor.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1), sync.Transfer.GetNonce())
}

func (s *ApplyTransferTestSuite) setUserStatesInTree() {
	senderStateID := senderState.PubKeyID
	receiverStateID := receiverState.PubKeyID

	_, err := s.storage.StateTree.Set(senderStateID, &senderState)
	s.NoError(err)
	_, err = s.storage.StateTree.Set(receiverStateID, &receiverState)
	s.NoError(err)
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
