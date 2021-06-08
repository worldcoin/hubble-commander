package storage

import (
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
			Hash:                 utils.RandomHash(),
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
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)

	err = s.storage.AddAccountIfNotExists(&account2)
	s.NoError(err)
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

func (s *TransferTestSuite) TestGetTransferWithBatchHash() {
	batch := &models.Batch{
		Type:            txtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		Number:          models.MakeUint256(1),
	}
	batchID, err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := commitment
	commitmentInBatch.IncludedInBatch = batchID
	commitmentID, err := s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInBatch := transfer
	transferInBatch.IncludedInCommitment = commitmentID
	err = s.storage.AddTransfer(&transferInBatch)
	s.NoError(err)

	expected := models.TransferWithBatchHash{
		Transfer:  transferInBatch,
		BatchHash: batch.Hash,
	}
	res, err := s.storage.GetTransferWithBatchHash(transferInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransferTestSuite) TestBatchAddTransfer() {
	txs := make([]models.Transfer, 2)
	txs[0] = transfer
	txs[0].Hash = utils.RandomHash()
	txs[1] = transfer
	txs[1].Hash = utils.RandomHash()

	err := s.storage.BatchAddTransfer(txs)
	s.NoError(err)

	transfer, err := s.storage.GetTransfer(txs[0].Hash)
	s.NoError(err)
	s.Equal(txs[0], *transfer)
	transfer, err = s.storage.GetTransfer(txs[1].Hash)
	s.NoError(err)
	s.Equal(txs[1], *transfer)
}

func (s *TransferTestSuite) TestBatchAddTransfer_NoTransfers() {
	err := s.storage.BatchAddTransfer([]models.Transfer{})
	s.Equal(ErrNoRowsAffected, err)
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
	transfer.Nonce = models.MakeUint256(1)
	transfer.Hash = utils.RandomHash()
	transfer2 := transfer
	transfer2.Nonce = models.MakeUint256(4)
	transfer2.Hash = utils.RandomHash()
	transfer3 := transfer
	transfer3.Nonce = models.MakeUint256(7)
	transfer3.Hash = utils.RandomHash()
	transfer4 := transfer
	transfer4.Nonce = models.MakeUint256(5)
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
		{PubKeyID: 1}, // StateID: 0
		{PubKeyID: 2}, // StateID: 1
		{PubKeyID: 1}, // StateID: 2
		{PubKeyID: 3}, // StateID: 3
		{PubKeyID: 2}, // StateID: 4
	}

	for i := range userStates {
		err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	batchHash, commitmentID := s.addBatchAndCommitment()
	transfers := make([]models.TransferWithBatchHash, 5)

	transfers[0].Transfer = transfer
	transfers[0].Hash = utils.RandomHash()
	transfers[0].FromStateID = 0
	transfers[0].ToStateID = 1
	transfers[0].IncludedInCommitment = &commitmentID
	transfers[0].BatchHash = &batchHash

	transfers[1].Transfer = transfer
	transfers[1].Hash = utils.RandomHash()
	transfers[1].FromStateID = 1
	transfers[1].ToStateID = 4
	transfers[1].IncludedInCommitment = &commitmentID
	transfers[1].BatchHash = &batchHash

	transfers[2].Transfer = transfer
	transfers[2].Hash = utils.RandomHash()
	transfers[2].FromStateID = 2
	transfers[2].ToStateID = 1

	transfers[3].Transfer = transfer
	transfers[3].Hash = utils.RandomHash()
	transfers[3].FromStateID = 3
	transfers[3].ToStateID = 1
	transfers[3].IncludedInCommitment = &commitmentID
	transfers[3].BatchHash = &batchHash

	transfers[4].Transfer = transfer
	transfers[4].Hash = utils.RandomHash()
	transfers[4].FromStateID = 1
	transfers[4].ToStateID = 2

	for i := range transfers {
		err := s.storage.AddTransfer(&transfers[i].Transfer)
		s.NoError(err)
	}

	userTransactions, err := s.storage.GetTransfersByPublicKey(&models.PublicKey{1, 2, 3})
	s.NoError(err)
	s.Len(userTransactions, 4)
	s.Contains(userTransactions, transfers[0])
	s.Contains(userTransactions, transfers[2])
	s.Contains(userTransactions, transfers[3])
	s.Contains(userTransactions, transfers[4])
}

func (s *TransferTestSuite) TestGetUserTransfersByPublicKey_NoTransfers() {
	userTransfers, err := s.storage.GetTransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(userTransfers, 0)
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

func (s *TransferTestSuite) addBatchAndCommitment() (batchHash common.Hash, commitmentID int32) {
	batch := &models.Batch{
		Type:            txtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		Number:          models.MakeUint256(1),
	}
	batchID, err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := commitment
	commitmentInBatch.IncludedInBatch = batchID
	id, err := s.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)

	return *batch.Hash, *id
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
