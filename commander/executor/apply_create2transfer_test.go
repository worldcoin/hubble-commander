package executor

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
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
	storage             *st.Storage
	teardown            func() error
	tree                *st.StateTree
	transactionExecutor *TransactionExecutor
	client              *eth.TestClient
}

func (s *ApplyCreate2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyCreate2TransferTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.tree = st.NewStateTree(s.storage)
	s.client, err = eth.NewTestClient()
	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, nil, TransactionExecutorOpts{})
	s.NoError(err)

	accounts := []models.Account{
		{
			PubKeyID:  0,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{3, 4, 5},
		},
	}
	for i := range accounts {
		err = s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	_, err = s.tree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  feeReceiverTokenID,
		Balance:  models.MakeUint256(10000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
	_, err = s.tree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  feeReceiverTokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_InsertsNewEmptyStateLeaf() {
	c2T := create2Transfer
	transferError, appError := s.transactionExecutor.ApplyCreate2Transfer(&c2T, 2, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(transferError)

	leaf, err := s.storage.GetStateLeaf(2)
	s.NoError(err)
	s.NotNil(leaf)
	s.Equal(models.MakeUint256(0), leaf.Nonce)
	s.Equal(feeReceiverTokenID, leaf.TokenID)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_AppliesTransfer() {
	c2T := create2Transfer
	transferError, appError := s.transactionExecutor.ApplyCreate2Transfer(&c2T, 2, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(transferError)

	receiverLeaf, err := s.storage.GetStateLeaf(2)
	s.NoError(err)
	senderLeaf, err := s.storage.GetStateLeaf(c2T.FromStateID)
	s.NoError(err)

	s.Equal(uint64(8900), senderLeaf.Balance.Uint64())
	s.Equal(uint64(1000), receiverLeaf.Balance.Uint64())
}

// TODO-AFS rework this test and add tests for ApplyCreate2TransferForSync
func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_TransferWithStateID() {
	c2T := create2Transfer
	c2T.ToStateID = ref.Uint32(5)
	transferError, appError := s.transactionExecutor.ApplyCreate2Transfer(&c2T, 2, feeReceiverTokenID)
	s.NoError(appError)
	s.NoError(transferError)
	s.Equal(uint32(5), *c2T.ToStateID)

	receiverLeaf, err := s.storage.GetStateLeaf(*c2T.ToStateID)
	s.NoError(err)
	senderLeaf, err := s.storage.GetStateLeaf(c2T.FromStateID)
	s.NoError(err)

	s.Equal(uint64(8900), senderLeaf.Balance.Uint64())
	s.Equal(uint64(1000), receiverLeaf.Balance.Uint64())
}

func (s *ApplyCreate2TransfersTestSuite) TestApplyCreate2Transfer_InvalidTransfer() {
	transfers := generateInvalidCreate2Transfers(1, &s.publicKey)
	transfers[0].Amount = models.MakeUint256(500)

	transferErr, appErr := s.transactionExecutor.ApplyCreate2Transfer(&transfers[0], 1, *models.NewUint256(1))
	s.Error(transferErr)
	s.NoError(appErr)
}

func TestApplyCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransferTestSuite))
}
