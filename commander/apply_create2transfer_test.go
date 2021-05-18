package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	feeReceiverTokenIndex = models.MakeUint256(5)
)

type ApplyCreate2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage     *st.Storage
	teardown    func() error
	tree        *st.StateTree
	client      *eth.TestClient
	events      chan *accountregistry.AccountRegistryPubkeyRegistered
	unsubscribe func()
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

	err = s.tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: feeReceiverTokenIndex,
		Balance:    models.MakeUint256(10000),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
	err = s.tree.Set(1, &models.UserState{
		PubKeyID:   1,
		TokenIndex: feeReceiverTokenIndex,
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	s.events, s.unsubscribe, err = s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TearDownTest() {
	s.unsubscribe()
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_InsertsNewEmptyStateLeaf() {
	_, transferError, appError := ApplyCreate2Transfer(s.storage, s.client.Client, s.events, &create2Transfer, feeReceiverTokenIndex)
	s.NoError(appError)
	s.NoError(transferError)

	leaf, err := s.tree.Leaf(2)
	s.NoError(err)
	s.NotNil(leaf)
	s.Equal(models.MakeUint256(0), leaf.Nonce)
	s.Equal(feeReceiverTokenIndex, leaf.TokenIndex)
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_ApplyTransfer() {
	_, transferError, appError := ApplyCreate2Transfer(s.storage, s.client.Client, s.events, &create2Transfer, feeReceiverTokenIndex)
	s.NoError(appError)
	s.NoError(transferError)

	receiverLeaf, err := s.tree.Leaf(2)
	s.NoError(err)
	senderLeaf, err := s.tree.Leaf(create2Transfer.FromStateID)
	s.NoError(err)

	s.Equal(int64(8900), senderLeaf.Balance.Int64())
	s.Equal(int64(1000), receiverLeaf.Balance.Int64())
}

func (s *ApplyCreate2TransferTestSuite) TestApplyCreate2Transfer_ReturnsCorrectPubKeyID() {
	addedPubKeyID, transferError, appError :=
		ApplyCreate2Transfer(s.storage, s.client.Client, s.events, &create2Transfer, feeReceiverTokenIndex)
	s.NoError(appError)
	s.NoError(transferError)
	s.Equal(uint32(2), *addedPubKeyID)
}

func (s *ApplyCreate2TransferTestSuite) TestGetPubKeyID_AccountNotExists() {
	transfer := create2Transfer
	transfer.ToPublicKey = models.PublicKey{10, 11, 12}
	pubKeyID, err := getOrRegisterPubKeyID(s.storage, s.client.Client, s.events, &transfer, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(0), *pubKeyID)
}

func (s *ApplyCreate2TransferTestSuite) TestGetPubKeyID_AccountForTokenIndexNotExists() {
	pubKeyID, err := getOrRegisterPubKeyID(s.storage, s.client.Client, s.events, &create2Transfer, models.MakeUint256(1))
	s.NoError(err)
	s.Equal(uint32(2), *pubKeyID)
}

func TestApplyCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyCreate2TransferTestSuite))
}
