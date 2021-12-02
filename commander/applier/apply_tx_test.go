package applier

import (
	"testing"

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

type ApplyTxTestSuite struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	applier      *Applier
	transfer     models.Transfer
	receiverLeaf models.StateLeaf
}

func (s *ApplyTxTestSuite) SetupSuite() {
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

func (s *ApplyTxTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage)

	s.receiverLeaf = models.StateLeaf{
		StateID:   receiverState.PubKeyID,
		UserState: receiverState,
	}
}

func (s *ApplyTxTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyTxTestSuite) TestCalculateStateAfterTx_UpdatesStates() {
	newSenderState, newReceiverState, err := calculateStateAfterTx(
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

func (s *ApplyTxTestSuite) TestCalculateStateAfterTx_ValidatesTokenAmount() {
	invalidTransfer := s.transfer
	invalidTransfer.Amount = models.MakeUint256(0)
	_, _, err := calculateStateAfterTx(
		senderState,
		receiverState,
		&invalidTransfer,
	)
	s.ErrorIs(err, ErrInvalidTokenAmount)
}

func (s *ApplyTxTestSuite) TestCalculateStateAfterTx_ValidatesBalance() {
	transferAboveBalance := s.transfer
	transferAboveBalance.Amount = models.MakeUint256(410)

	_, _, err := calculateStateAfterTx(senderState, receiverState, &transferAboveBalance)
	s.ErrorIs(err, ErrBalanceTooLow)
}

func (s *ApplyTxTestSuite) TestCalculateStateAfterTx_ReturnsCorrectLeavesInCaseOfSelfTransfer() {
	selfTransfer := s.transfer
	selfTransfer.ToStateID = selfTransfer.FromStateID
	newSenderState, newReceiverState, err := calculateStateAfterTx(
		senderState,
		senderState,
		&selfTransfer,
	)
	s.NoError(err)

	s.Equal(models.MakeUint256(1), newSenderState.Nonce)
	s.Equal(models.MakeUint256(290), newSenderState.Balance)

	s.Equal(models.MakeUint256(1), newReceiverState.Nonce)
	s.Equal(models.MakeUint256(390), newReceiverState.Balance)

	s.NotEqual(&newSenderState, &senderState)
	s.NotEqual(&newReceiverState, &senderState)
}

func (s *ApplyTxTestSuite) TestApplyTx_ValidatesSenderTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	txError, appError := s.applier.ApplyTx(&s.transfer, &s.receiverLeaf, models.MakeUint256(3))
	s.NoError(txError)
	s.ErrorIs(appError, ErrInvalidSenderTokenID)
}

func (s *ApplyTxTestSuite) TestApplyTx_ValidatesReceiverTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	receiverWithChangedToken := s.receiverLeaf
	receiverWithChangedToken.TokenID = models.MakeUint256(2)

	txError, appError := s.applier.ApplyTx(&s.transfer, &receiverWithChangedToken, models.MakeUint256(1))
	s.NoError(txError)
	s.ErrorIs(appError, ErrInvalidReceiverTokenID)
}

func (s *ApplyTxTestSuite) TestApplyTx_ValidatesNonce() {
	transferWithBadNonce := s.transfer
	transferWithBadNonce.Nonce = models.MakeUint256(1)
	setUserStatesInTree(s.Assertions, s.storage)

	txError, appError := s.applier.ApplyTx(&transferWithBadNonce, &s.receiverLeaf, models.MakeUint256(1))
	s.ErrorIs(txError, ErrNonceTooHigh)
	s.NoError(appError)
}

func (s *ApplyTxTestSuite) TestApplyTx_UpdatesStatesCorrectly() {
	setUserStatesInTree(s.Assertions, s.storage)

	txError, appError := s.applier.ApplyTx(&s.transfer, &s.receiverLeaf, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func setUserStatesInTree(s *require.Assertions, storage *st.TestStorage) {
	senderStateID := senderState.PubKeyID
	receiverStateID := receiverState.PubKeyID

	_, err := storage.StateTree.Set(senderStateID, &senderState)
	s.NoError(err)
	_, err = storage.StateTree.Set(receiverStateID, &receiverState)
	s.NoError(err)
}

func TestApplyTxTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTxTestSuite))
}
