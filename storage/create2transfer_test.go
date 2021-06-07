package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	create2Transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			TxType:               txtype.Create2Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(1000),
			Fee:                  models.MakeUint256(100),
			Nonce:                models.MakeUint256(0),
			Signature:            models.MakeRandomSignature(),
			IncludedInCommitment: nil,
		},
		ToStateID:   ref.Uint32(2),
		ToPublicKey: account2.PublicKey,
	}
)

type Create2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	tree    *StateTree
}

func (s *Create2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)

	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TestAddCreate2Transfer_AddAndRetrieve() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	res, err := s.storage.GetCreate2Transfer(create2Transfer.Hash)
	s.NoError(err)

	s.Equal(create2Transfer, *res)
}

func (s *Create2TransferTestSuite) TestBatchAddCreate2Transfer() {
	txs := make([]models.Create2Transfer, 2)
	txs[0] = create2Transfer
	txs[0].Hash = utils.RandomHash()
	txs[1] = create2Transfer
	txs[1].Hash = utils.RandomHash()

	err := s.storage.BatchAddCreate2Transfer(txs)
	s.NoError(err)

	transfer, err := s.storage.GetCreate2Transfer(txs[0].Hash)
	s.NoError(err)
	s.Equal(txs[0], *transfer)
	transfer, err = s.storage.GetCreate2Transfer(txs[1].Hash)
	s.NoError(err)
	s.Equal(txs[1], *transfer)
}

func (s *Create2TransferTestSuite) TestBatchAddCreate2Transfer_NoTransfers() {
	err := s.storage.BatchAddCreate2Transfer([]models.Create2Transfer{})
	s.Equal(ErrNoRowsAffected, err)
}

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_NonExistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetCreate2Transfer(hash)
	s.Equal(NewNotFoundError("transaction"), err)
	s.Nil(res)
}

func (s *Create2TransferTestSuite) TestGetPendingCreate2Transfers() {
	commitment := &models.Commitment{}
	id, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	create2Transfer2 := create2Transfer
	create2Transfer2.Hash = utils.RandomHash()
	create2Transfer3 := create2Transfer
	create2Transfer3.Hash = utils.RandomHash()
	create2Transfer3.IncludedInCommitment = id
	create2Transfer4 := create2Transfer
	create2Transfer4.Hash = utils.RandomHash()
	create2Transfer4.ErrorMessage = ref.String("A very boring error message")

	create2Transfers := []*models.Create2Transfer{&create2Transfer, &create2Transfer2, &create2Transfer3, &create2Transfer4}

	for _, create2Transfer := range create2Transfers {
		err = s.storage.AddCreate2Transfer(create2Transfer)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)

	s.Equal([]models.Create2Transfer{create2Transfer, create2Transfer2}, res)
}

func (s *Create2TransferTestSuite) TestGetPendingCreate2Transfers_OrdersTransfersByNonceAscending() {
	create2Transfer.TransactionBase.Nonce = models.MakeUint256(1)
	create2Transfer.Hash = utils.RandomHash()
	create2Transfer2 := create2Transfer
	create2Transfer2.TransactionBase.Nonce = models.MakeUint256(4)
	create2Transfer2.Hash = utils.RandomHash()
	create2Transfer3 := create2Transfer
	create2Transfer3.TransactionBase.Nonce = models.MakeUint256(7)
	create2Transfer3.Hash = utils.RandomHash()
	create2Transfer4 := create2Transfer
	create2Transfer4.TransactionBase.Nonce = models.MakeUint256(5)
	create2Transfer4.Hash = utils.RandomHash()

	for _, transfer := range []*models.Create2Transfer{&create2Transfer, &create2Transfer2, &create2Transfer3, &create2Transfer4} {
		err := s.storage.AddCreate2Transfer(transfer)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)

	s.Equal([]models.Create2Transfer{create2Transfer, create2Transfer2, create2Transfer4, create2Transfer3}, res)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByPublicKey() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	err = s.tree.Set(1, &models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(400),
	})
	s.NoError(err)

	transfers, err := s.storage.GetCreate2TransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(transfers, 1)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByPublicKey_NoCreate2Transfers() {
	transfers, err := s.storage.GetCreate2TransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByCommitmentID() {
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	transfer1 := create2Transfer
	transfer1.IncludedInCommitment = commitmentID

	err = s.storage.AddCreate2Transfer(&transfer1)
	s.NoError(err)

	commitments, err := s.storage.GetCreate2TransfersByCommitmentID(*commitmentID)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByCommitmentID_NoCreate2Transfers() {
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	commitments, err := s.storage.GetCreate2TransfersByCommitmentID(*commitmentID)
	s.NoError(err)
	s.Len(commitments, 0)
}

func (s *Create2TransferTestSuite) TestSetCreate2TransferToStateID() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	toStateID := uint32(10)
	err = s.storage.SetCreate2TransferToStateID(create2Transfer.Hash, toStateID)
	s.NoError(err)

	c2t, err := s.storage.GetCreate2Transfer(create2Transfer.Hash)
	s.NoError(err)
	s.Equal(toStateID, *c2t.ToStateID)
}

func (s *Create2TransferTestSuite) TestSetCreate2TransferToStateID_NoCreate2Transfer() {
	err := s.storage.SetCreate2TransferToStateID(create2Transfer.Hash, 10)
	s.Equal(err, ErrNoRowsAffected)
}

func TestCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferTestSuite))
}
