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
	transfer = models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.BigToHash(big.NewInt(1234)),
			TxType:               txtype.Transfer,
			FromStateID:          1,
			Amount:               models.MakeUint256(1000),
			Fee:                  models.MakeUint256(100),
			Nonce:                models.MakeUint256(0),
			Signature:            models.MakeRandomSignature(),
			IncludedInCommitment: nil,
		},
		ToStateID: 2,
	}
)

type TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	tree    *StateTree
}

func (s *TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)
}

func (s *TransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransferTestSuite) TestAddTransfer_AddAndRetrieve() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)

	s.Equal(transfer, *res)
}

func (s *TransferTestSuite) TestGetTransfer_NonExistentTransfer() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetTransfer(hash)
	s.Equal(NewNotFoundError("transaction"), err)
	s.Nil(res)
}

func (s *TransferTestSuite) TestGetPendingTransfers() {
	commitment := &models.Commitment{}
	id, err := s.storage.AddCommitment(commitment)
	s.NoError(err)

	transfer2 := transfer
	transfer2.Hash = utils.RandomHash()
	transfer3 := transfer
	transfer3.Hash = utils.RandomHash()
	transfer3.IncludedInCommitment = id
	transfer4 := transfer
	transfer4.Hash = utils.RandomHash()
	transfer4.ErrorMessage = ref.String("A very boring error message")

	for _, transfer := range []*models.Transfer{&transfer, &transfer2, &transfer3, &transfer4} {
		err = s.storage.AddTransfer(transfer)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingTransfers()
	s.NoError(err)

	s.Equal([]models.Transfer{transfer, transfer2}, res)
}

func (s *TransferTestSuite) TestGetPendingTransfers_OrdersTransfersByNonceAscending() {
	transfer.TransactionBase.Nonce = models.MakeUint256(1)
	transfer.Hash = utils.RandomHash()
	transfer2 := transfer
	transfer2.TransactionBase.Nonce = models.MakeUint256(4)
	transfer2.Hash = utils.RandomHash()
	transfer3 := transfer
	transfer3.TransactionBase.Nonce = models.MakeUint256(7)
	transfer3.Hash = utils.RandomHash()
	transfer4 := transfer
	transfer4.TransactionBase.Nonce = models.MakeUint256(5)
	transfer4.Hash = utils.RandomHash()

	for _, transfer := range []*models.Transfer{&transfer, &transfer2, &transfer3, &transfer4} {
		err := s.storage.AddTransfer(transfer)
		s.NoError(err)
	}

	res, err := s.storage.GetPendingTransfers()
	s.NoError(err)

	s.Equal([]models.Transfer{transfer, transfer2, transfer4, transfer3}, res)
}

func (s *TransferTestSuite) TestGetUserTransfers() {
	transfer1 := transfer
	transfer1.Hash = utils.RandomHash()
	transfer1.FromStateID = 1
	transfer2 := transfer
	transfer2.Hash = utils.RandomHash()
	transfer2.FromStateID = 2
	transfer3 := transfer
	transfer3.Hash = utils.RandomHash()
	transfer3.FromStateID = 1

	err := s.storage.AddTransfer(&transfer1)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer2)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer3)
	s.NoError(err)

	userTransactions, err := s.storage.GetUserTransfers(models.MakeUint256(1))
	s.NoError(err)

	s.Len(userTransactions, 2)
	s.Contains(userTransactions, transfer1)
	s.Contains(userTransactions, transfer3)
}

func (s *TransferTestSuite) TestGetUserTransfers_NoTransfers() {
	userTransactions, err := s.storage.GetUserTransfers(models.MakeUint256(1))

	s.NoError(err)
	s.Len(userTransactions, 0)
}

func (s *TransferTestSuite) TestGetTransfersByPublicKey() {
	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   2,
			TokenIndex: models.MakeUint256(2),
			Balance:    models.MakeUint256(500),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(25),
			Balance:    models.MakeUint256(1),
			Nonce:      models.MakeUint256(73),
		},
		{
			PubKeyID:   3,
			TokenIndex: models.MakeUint256(30),
			Balance:    models.MakeUint256(50),
			Nonce:      models.MakeUint256(71),
		},
	}

	for i := range userStates {
		err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	transfer1 := transfer
	transfer1.Hash = utils.RandomHash()
	transfer1.FromStateID = 0
	transfer2 := transfer
	transfer2.Hash = utils.RandomHash()
	transfer2.FromStateID = 1
	transfer3 := transfer
	transfer3.Hash = utils.RandomHash()
	transfer3.FromStateID = 2
	transfer4 := transfer
	transfer4.Hash = utils.RandomHash()
	transfer4.FromStateID = 3

	err := s.storage.AddTransfer(&transfer1)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer2)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer3)
	s.NoError(err)
	err = s.storage.AddTransfer(&transfer4)
	s.NoError(err)

	userTransactions, err := s.storage.GetTransfersByPublicKey(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(userTransactions, 3)
	s.Contains(userTransactions, transfer1)
	s.Contains(userTransactions, transfer3)
	s.Contains(userTransactions, transfer4)
}

func (s *TransferTestSuite) TestGetUserTransfersByPublicKey_NoTransfers() {
	userTransfers, err := s.storage.GetTransfersByPublicKey(&models.PublicKey{1, 2, 3})

	s.NoError(err)
	s.Len(userTransfers, 0)
	s.NotNil(userTransfers)
}

func (s *TransferTestSuite) TestGetTransfersByCommitmentID() {
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	transfer1 := transfer
	transfer1.IncludedInCommitment = commitmentID

	err = s.storage.AddTransfer(&transfer1)
	s.NoError(err)

	commitments, err := s.storage.GetTransfersByCommitmentID(*commitmentID)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *TransferTestSuite) TestGetTransfersByCommitmentID_NoTransfers() {
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	commitments, err := s.storage.GetTransfersByCommitmentID(*commitmentID)
	s.NoError(err)
	s.Len(commitments, 0)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
