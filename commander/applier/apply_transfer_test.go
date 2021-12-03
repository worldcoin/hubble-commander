package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	applier      *Applier
	transfer     models.Transfer
	receiverLeaf models.StateLeaf
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

	s.receiverLeaf = models.StateLeaf{
		StateID:   receiverState.PubKeyID,
		UserState: receiverState,
	}
}

func (s *ApplyTransferTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage, nil)
}

func (s *ApplyTransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, txError, appError := s.applier.ApplyTransfer(&s.transfer, s.receiverLeaf.TokenID)
	s.NoError(txError)
	s.NoError(appError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_ValidatesReceiverStateID() {
	setUserStatesInTree(s.Assertions, s.storage)

	transfer := s.transfer
	transfer.ToStateID = 10

	_, txError, appError := s.applier.ApplyTransfer(&transfer, s.receiverLeaf.TokenID)
	s.NoError(appError)
	s.ErrorIs(txError, ErrNonexistentReceiver)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsSenderProofForCalculateStateAfterTransferValidations() {
	setUserStatesInTree(s.Assertions, s.storage)

	bigTransfer := s.transfer
	bigTransfer.Amount = models.MakeUint256(1_000_000)

	synced, txError, appError := s.applier.ApplyTransferForSync(&bigTransfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.ErrorIs(txError, ErrBalanceTooLow)
	s.NoError(appError)

	s.Equal(&bigTransfer, synced.Tx.ToTransfer())
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesSenderTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	synced, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(3))
	s.NotNil(synced)
	s.ErrorIs(txError, ErrInvalidSenderTokenID)
	s.NoError(appError)

	s.Equal(&s.transfer, synced.Tx.ToTransfer())
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesReceiverTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	receiverWithChangedToken := receiverState
	receiverWithChangedToken.TokenID = models.MakeUint256(2)
	_, err := s.storage.StateTree.Set(2, &receiverWithChangedToken)
	s.NoError(err)

	synced, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.ErrorIs(txError, ErrInvalidReceiverTokenID)
	s.NoError(appError)

	s.Equal(&s.transfer, synced.Tx.ToTransfer())
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
	s.Equal(receiverWithChangedToken, *synced.ReceiverStateProof.UserState)
	s.Len(synced.ReceiverStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsTransferWithUpdatedNonce() {
	setUserStatesInTree(s.Assertions, s.storage)

	transferWithModifiedNonce := s.transfer
	transferWithModifiedNonce.Nonce = models.MakeUint256(1234)

	synced, txError, appError := s.applier.ApplyTransferForSync(&transferWithModifiedNonce, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	s.Equal(models.MakeUint256(1234), transferWithModifiedNonce.Nonce)
	s.Equal(models.MakeUint256(0), synced.Tx.ToTransfer().GetNonce())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_UpdatesStatesCorrectly() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsProofs() {
	setUserStatesInTree(s.Assertions, s.storage)

	sync, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	s.Equal(senderState, *sync.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
	s.Equal(receiverState, *sync.ReceiverStateProof.UserState)
	s.Len(sync.ReceiverStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesNonexistentSenderState() {
	setUserStatesInTree(s.Assertions, s.storage)

	senderLeaf, err := s.storage.StateTree.LeafOrEmpty(10)
	s.NoError(err)

	transfer := s.transfer
	transfer.FromStateID = senderLeaf.StateID

	sync, txError, appError := s.applier.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.ErrorIs(txError, ErrBalanceTooLow)
	s.Equal(senderLeaf.UserState, *sync.Proofs.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_AllowsNonexistentReceiverState() {
	_, err := s.storage.StateTree.Set(s.transfer.FromStateID, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(400),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	expectedReceiverLeaf, err := st.NewStateLeaf(2, &models.UserState{
		Balance: models.MakeUint256(100),
	})
	s.NoError(err)

	_, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(0))
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(s.transfer.FromStateID)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(s.transfer.ToStateID)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(expectedReceiverLeaf, receiverLeaf)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_SetsNonce() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	sync, txError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(txError)

	s.Equal(models.MakeUint256(1), sync.Tx.ToTransfer().GetNonce())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_AllowTheSameFromTo() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}

	_, err := s.storage.StateTree.Set(transfer.FromStateID, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(400),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, txError, appError := s.applier.ApplyTransferForSync(&transfer, models.MakeUint256(0))
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(transfer.FromStateID)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(transfer.ToStateID)
	s.NoError(err)

	s.EqualValues(390, senderLeaf.Balance.Uint64())
	s.EqualValues(1, senderLeaf.Nonce.Uint64())
	s.Equal(receiverLeaf, senderLeaf)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_TheSameFromToLowBalance() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(20),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 1,
	}

	_, err := s.storage.StateTree.Set(transfer.FromStateID, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(50),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, txError, appError := s.applier.ApplyTransferForSync(&transfer, models.MakeUint256(0))
	s.ErrorIs(txError, ErrBalanceTooLow)
	s.NoError(appError)

	senderLeaf, err := s.storage.StateTree.Leaf(transfer.FromStateID)
	s.NoError(err)

	s.EqualValues(50, senderLeaf.Balance.Uint64())
	s.EqualValues(0, senderLeaf.Nonce.Uint64())
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
