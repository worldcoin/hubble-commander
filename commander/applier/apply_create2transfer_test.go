package applier

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	create2Transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			FromStateID: 0,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
		},
		ToPublicKey: models.PublicKey{3, 4, 5},
	}
	feeReceiverTokenID = models.MakeUint256(5)
)

type ApplyCreate2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	client  *eth.TestClient
	applier *Applier
}

func (s *ApplyCreate2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransferTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.applier = NewApplier(s.storage.Storage)

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  feeReceiverTokenID,
		Balance:  models.MakeUint256(10000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  feeReceiverTokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_GetsNextAvailableStateIDAndInsertsNewUserState() {
	_, txError, appError := s.applier.ApplyCreate2Transfer(&create2Transfer, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	leaf, err := s.storage.StateTree.Leaf(2)
	s.NoError(err)
	s.NotNil(leaf)
	s.EqualValues(st.AccountBatchOffset, leaf.PubKeyID)
	s.Equal(feeReceiverTokenID, leaf.TokenID)
	s.Equal(models.MakeUint256(0), leaf.Nonce)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_SetsCorrectToStateIDInReturnedTransfer() {
	applyResult, txError, appError := s.applier.ApplyCreate2Transfer(&create2Transfer, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	s.Equal(ref.Uint32(2), applyResult.AppliedTx().GetToStateID())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_AppliesTransfer() {
	applyResult, txError, appError := s.applier.ApplyCreate2Transfer(&create2Transfer, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(applyResult.AppliedTx().GetFromStateID())
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(*applyResult.AppliedTx().GetToStateID())
	s.NoError(err)

	s.Equal(uint64(8900), senderLeaf.Balance.Uint64())
	s.Equal(uint64(1000), receiverLeaf.Balance.Uint64())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_InvalidTransfer() {
	transfers := testutils.GenerateInvalidCreate2Transfers(1, 1)
	transfers[0].Amount = models.MakeUint256(500)

	_, txError, appErr := s.applier.ApplyCreate2Transfer(&transfers[0], feeReceiverTokenID)
	s.Error(txError)
	s.NoError(appErr)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_ReturnsPendingAccountWhenPublicKeyIsUnknown() {
	applyResult, txError, appError := s.applier.ApplyCreate2Transfer(&create2Transfer, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	expectedAccount := models.AccountLeaf{
		PubKeyID:  st.AccountBatchOffset,
		PublicKey: create2Transfer.ToPublicKey,
	}
	s.NotNil(applyResult.PendingAccount())
	s.Equal(expectedAccount, *applyResult.PendingAccount())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_DoesNotReturnPendingAccountWhenPublicKeyIsRegistered() {
	err := s.applier.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  2,
		PublicKey: create2Transfer.ToPublicKey,
	})
	s.NoError(err)

	applyResult, txError, appError := s.applier.ApplyCreate2Transfer(&create2Transfer, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)
	s.Nil(applyResult.PendingAccount())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2TransferForSync_ReturnsErrorOnNilToStateID() {
	_, txError, appError := s.applier.ApplyCreate2TransferForSync(&create2Transfer, uint32(2), feeReceiverTokenID)
	s.NoError(txError)
	s.ErrorIs(appError, ErrNilReceiverStateID)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2TransferForSync_InsertsNewUserStateAtReceiverStateID() {
	pubKeyID := uint32(2)
	c2T := create2Transfer
	c2T.ToStateID = ref.Uint32(5)
	_, txError, appError := s.applier.ApplyCreate2TransferForSync(&c2T, pubKeyID, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	leaf, err := s.storage.StateTree.Leaf(*c2T.ToStateID)
	s.NoError(err)

	s.NoError(err)
	s.NotNil(leaf)
	s.Equal(pubKeyID, leaf.PubKeyID)
	s.Equal(feeReceiverTokenID, leaf.TokenID)
	s.Equal(models.MakeUint256(0), leaf.Nonce)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2TransferForSync_AppliesTransfer() {
	c2T := create2Transfer
	c2T.ToStateID = ref.Uint32(5)
	_, txError, appError := s.applier.ApplyCreate2TransferForSync(&c2T, 2, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(txError)

	senderLeaf, err := s.storage.StateTree.Leaf(c2T.FromStateID)
	s.NoError(err)
	receiverLeaf, err := s.storage.StateTree.Leaf(*c2T.ToStateID)
	s.NoError(err)

	s.Equal(uint64(8900), senderLeaf.Balance.Uint64())
	s.Equal(uint64(1000), receiverLeaf.Balance.Uint64())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2TransferForSync_ValidatesNonexistentSenderState() {
	senderLeaf, err := s.storage.StateTree.LeafOrEmpty(10)
	s.NoError(err)

	transfer := create2Transfer
	transfer.ToStateID = ref.Uint32(5)
	transfer.FromStateID = senderLeaf.StateID

	sync, txError, appError := s.applier.ApplyCreate2TransferForSync(&transfer, 2, feeReceiverTokenID)
	s.NoError(appError)
	s.ErrorIs(txError, ErrBalanceTooLow)
	s.Equal(senderLeaf.UserState, *sync.Proofs.SenderStateProof.UserState)
	s.Len(sync.SenderStateProof.Witness, st.StateTreeDepth)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2TransferForSync_InvalidTransfer() {
	transfers := testutils.GenerateValidCreate2Transfers(1)
	invalidC2T := transfers[0]
	invalidC2T.Amount = models.MakeUint256(1_000_000)
	invalidC2T.ToStateID = ref.Uint32(5)

	_, txError, appErr := s.applier.ApplyCreate2TransferForSync(&invalidC2T, 1, feeReceiverTokenID)
	s.Error(txError)
	s.NoError(appErr)
}

func (s *ApplyCreate2TransferTestSuite) TestGetPubKeyID_WhenPubKeyIsNotRegistered_PredictsPubKeyID() {
	pubKeyID, isPending, err := s.applier.getPubKeyID(&create2Transfer.ToPublicKey)
	s.NoError(err)
	s.True(isPending)
	s.Equal(uint32(st.AccountBatchOffset), *pubKeyID)
}

func (s *ApplyCreate2TransferTestSuite) TestGetPubKeyID_WhenPubKeyIsRegistered_ReturnsFirstPubKeyID() {
	for i := 1; i <= 3; i++ {
		err := s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: models.PublicKey{1, 2, 3},
		})
		s.NoError(err)
	}

	pubKeyID, isPending, err := s.applier.getPubKeyID(&models.PublicKey{1, 2, 3})
	s.NoError(err)
	s.False(isPending)
	s.Equal(uint32(1), *pubKeyID)
}

func TestApplyCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransferTestSuite))
}
