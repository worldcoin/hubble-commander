package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
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
	s.applier = NewApplier(s.storage.Storage, &eth.Client{})
}

func (s *ApplyTransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyTransferTestSuite) TestApplyTransfer() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, transferError, appError := s.applier.ApplyTransfer(&s.transfer, s.receiverLeaf.TokenID)
	s.NoError(transferError)
	s.NoError(appError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(uint64(100), receiverLeaf.Balance.Uint64())
}

func (s *ApplyTransferTestSuite) TestApplyTransfer_NonexistentReceiver() {
	setUserStatesInTree(s.Assertions, s.storage)

	transfer := s.transfer
	transfer.ToStateID = 10

	_, transferError, appError := s.applier.ApplyTransfer(&transfer, s.receiverLeaf.TokenID)
	s.NoError(transferError)
	s.Equal(st.NewNotFoundError("state leaf"), appError)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ReturnsSenderProofForCalculateStateAfterTransferValidations() {
	setUserStatesInTree(s.Assertions, s.storage)

	bigTransfer := s.transfer
	bigTransfer.Amount = models.MakeUint256(1_000_000)

	synced, transferError, appError := s.applier.ApplyTransferForSync(&bigTransfer, models.MakeUint256(1))
	s.NotNil(synced)
	s.Equal(ErrBalanceTooLow, transferError)
	s.NoError(appError)

	s.Equal(&bigTransfer, synced.Transfer)
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesSenderTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	synced, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(3))
	s.NotNil(synced)
	s.Equal(ErrInvalidSenderTokenID, transferError)
	s.NoError(appError)

	s.Equal(&s.transfer, synced.Transfer)
	s.Equal(senderState, *synced.SenderStateProof.UserState)
	s.Len(synced.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesReceiverTokenID() {
	setUserStatesInTree(s.Assertions, s.storage)

	receiverWithChangedToken := receiverState
	receiverWithChangedToken.TokenID = models.MakeUint256(2)
	_, err := s.storage.StateTree.Set(2, &receiverWithChangedToken)
	s.NoError(err)

	synced, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
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
	setUserStatesInTree(s.Assertions, s.storage)

	transferWithModifiedNonce := s.transfer
	transferWithModifiedNonce.Nonce = models.MakeUint256(1234)

	synced, transferError, appError := s.applier.ApplyTransferForSync(&transferWithModifiedNonce, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1234), transferWithModifiedNonce.Nonce)
	s.Equal(models.MakeUint256(0), synced.Transfer.GetNonce())
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_UpdatesStatesCorrectly() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
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
	setUserStatesInTree(s.Assertions, s.storage)

	sync, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(senderState, *sync.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
	s.Equal(receiverState, *sync.ReceiverStateProof.UserState)
	s.Len(sync.ReceiverStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_ValidatesNotExistingSenderState() {
	setUserStatesInTree(s.Assertions, s.storage)

	senderLeaf, err := s.storage.StateTree.LeafOrEmpty(10)
	s.NoError(err)

	transfer := s.transfer
	transfer.FromStateID = senderLeaf.StateID

	sync, transferError, appError := s.applier.ApplyTransferForSync(&transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.ErrorIs(transferError, ErrBalanceTooLow)
	s.Equal(senderLeaf.UserState, *sync.Proofs.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_AllowsNotExistingReceiverState() {
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

	_, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(0))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.storage.StateTree.Leaf(s.transfer.FromStateID)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(s.transfer.ToStateID)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
	s.Equal(expectedReceiverLeaf, receiverLeaf)
}

func (s *ApplyTransferTestSuite) TestApplyTransferForSync_SetsNonce() {
	setUserStatesInTree(s.Assertions, s.storage)

	_, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	sync, transferError, appError := s.applier.ApplyTransferForSync(&s.transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	s.Equal(models.MakeUint256(1), sync.Transfer.GetNonce())
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
