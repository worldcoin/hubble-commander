package storage

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
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
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: 2,
	}
)

type TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *TransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransferTestSuite) TestAddTransfer_AddAndRetrieve() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	expected := transfer

	res, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransferTestSuite) TestAddTransfer_AddAndRetrieveIncludedTransfer() {
	includedTransfer := transfer
	includedTransfer.CommitmentID = &models.CommitmentID{
		BatchID:      models.MakeUint256(3),
		IndexInBatch: 1,
	}
	err := s.storage.AddTransfer(&includedTransfer)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Equal(includedTransfer, *res)
}

func (s *TransferTestSuite) TestGetTransfer_DifferentTxType() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	_, err = s.storage.GetTransfer(create2Transfer.Hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
}

func (s *TransferTestSuite) TestMarkTransfersAsIncluded() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transfer
		txs[i].Hash = utils.RandomHash()
		err := s.storage.AddTransfer(&txs[i])
		s.NoError(err)
	}

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 1,
	}
	err := s.storage.MarkTransfersAsIncluded(txs, &commitmentID)
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetTransfer(txs[i].Hash)
		s.NoError(err)
		s.Equal(commitmentID, *tx.CommitmentID)
	}
}

func (s *TransferTestSuite) TestGetTransferWithBatchDetails() {
	batch := &models.Batch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		SubmissionTime:  &models.Timestamp{Time: time.Unix(140, 0).UTC()},
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	transferInBatch := transfer
	transferInBatch.CommitmentID = &models.CommitmentID{
		BatchID: batch.ID,
	}
	err = s.storage.AddTransfer(&transferInBatch)
	s.NoError(err)

	expected := models.TransferWithBatchDetails{
		Transfer:  transferInBatch,
		BatchHash: batch.Hash,
		BatchTime: batch.SubmissionTime,
	}
	res, err := s.storage.GetTransferWithBatchDetails(transferInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransferTestSuite) TestGetTransferWithBatchDetails_WithoutBatch() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	expected := models.TransferWithBatchDetails{Transfer: transfer}

	res, err := s.storage.GetTransferWithBatchDetails(transfer.Hash)
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
	s.ErrorIs(err, ErrNoRowsAffected)
}

func (s *TransferTestSuite) TestGetTransfer_NonexistentTransfer() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetTransfer(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func (s *TransferTestSuite) TestGetPendingTransfers() {
	transfers := make([]models.Transfer, 4)
	for i := range transfers {
		transfers[i] = transfer
		transfers[i].Hash = utils.RandomHash()
	}
	transfers[2].CommitmentID = &models.CommitmentID{BatchID: models.MakeUint256(3)}
	transfers[3].ErrorMessage = ref.String("A very boring error message")

	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)

	res, err := s.storage.GetPendingTransfers()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, transfers[0])
	s.Contains(res, transfers[1])
}

func (s *TransferTestSuite) TestGetPendingTransfers_OrdersTransfersByNonceAndTxHashAscending() {
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
	transfer4.Hash = common.Hash{66, 66, 66, 66}
	transfer5 := transfer
	transfer5.Nonce = models.MakeUint256(5)
	transfer5.Hash = common.Hash{65, 65, 65, 65}

	transfers := []models.Transfer{
		transfer,
		transfer2,
		transfer3,
		transfer4,
		transfer5,
	}

	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)

	res, err := s.storage.GetPendingTransfers()
	s.NoError(err)

	s.Equal(models.TransferArray{transfer, transfer2, transfer5, transfer4, transfer3}, res)
}

func (s *TransferTestSuite) TestGetTransfersByPublicKey() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  4,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  5,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{PubKeyID: 3}, // StateID: 0
		{PubKeyID: 4}, // StateID: 1
		{PubKeyID: 3}, // StateID: 2
		{PubKeyID: 5}, // StateID: 3
		{PubKeyID: 4}, // StateID: 4
	}

	for i := range userStates {
		_, err := s.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}

	submissionTime := &models.Timestamp{Time: time.Unix(170, 0).UTC()}
	batch := s.addBatchAndCommitment()
	commitmentID := &models.CommitmentID{BatchID: batch.ID}
	transfers := make([]models.TransferWithBatchDetails, 5)

	transfers[0].Transfer = transfer
	transfers[0].Hash = utils.RandomHash()
	transfers[0].FromStateID = 0
	transfers[0].ToStateID = 1
	transfers[0].CommitmentID = commitmentID
	transfers[0].BatchHash = batch.Hash
	transfers[0].BatchTime = submissionTime

	transfers[1].Transfer = transfer
	transfers[1].Hash = utils.RandomHash()
	transfers[1].FromStateID = 1
	transfers[1].ToStateID = 4
	transfers[1].CommitmentID = commitmentID
	transfers[1].BatchHash = batch.Hash

	transfers[2].Transfer = transfer
	transfers[2].Hash = utils.RandomHash()
	transfers[2].FromStateID = 2
	transfers[2].ToStateID = 1

	transfers[3].Transfer = transfer
	transfers[3].Hash = utils.RandomHash()
	transfers[3].FromStateID = 3
	transfers[3].ToStateID = 1
	transfers[3].CommitmentID = commitmentID
	transfers[3].BatchHash = batch.Hash
	transfers[3].BatchTime = submissionTime

	transfers[4].Transfer = transfer
	transfers[4].Hash = utils.RandomHash()
	transfers[4].FromStateID = 1
	transfers[4].ToStateID = 2

	s.batchAddTransfers(transfers)

	userTransactions, err := s.storage.GetTransfersByPublicKey(&models.PublicKey{1, 2, 3})
	s.NoError(err)
	s.Len(userTransactions, 4)
	s.Contains(userTransactions, transfers[0])
	s.Contains(userTransactions, transfers[2])
	s.Contains(userTransactions, transfers[3])
	s.Contains(userTransactions, transfers[4])
}

func (s *TransferTestSuite) TestGetTransfersByPublicKey_NoTransfersUnregisteredAccount() {
	userTransfers, err := s.storage.GetTransfersByPublicKey(&models.PublicKey{9, 9, 9})
	s.NoError(err)
	s.Len(userTransfers, 0)
}

func (s *TransferTestSuite) TestGetTransfersByPublicKey_NoTransfersRegisteredAccount() {
	err := s.storage.AccountTree.SetSingle(&account2)
	s.NoError(err)
	userTransfers, err := s.storage.GetTransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(userTransfers, 0)
}

func (s *TransferTestSuite) TestGetTransfersByPublicKey_NoTransfersButSomeCreate2Transfers() {
	err := s.storage.AccountTree.SetSingle(&account2)
	s.NoError(err)

	err = s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	userTransfers, err := s.storage.GetTransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(userTransfers, 0)
}

func (s *TransferTestSuite) TestGetTransfersByCommitmentID() {
	transfer1 := transfer
	transfer1.CommitmentID = &txCommitment.ID

	err := s.storage.AddTransfer(&transfer1)
	s.NoError(err)

	transfers, err := s.storage.GetTransfersByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 1)
}

func (s *TransferTestSuite) TestGetTransfersByCommitmentID_NoTransactions() {
	transfers, err := s.storage.GetTransfersByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *TransferTestSuite) TestGetTransfersByCommitmentID_NoTransfersButSomeCreate2Transfers() {
	c2t := create2Transfer
	c2t.CommitmentID = &txCommitment.ID
	err := s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	transfers, err := s.storage.GetTransfersByCommitmentID(txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *TransferTestSuite) addBatchAndCommitment() *models.Batch {
	batch := &models.Batch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		SubmissionTime:  &models.Timestamp{Time: time.Unix(170, 0).UTC()},
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := txCommitment
	commitmentInBatch.ID.BatchID = batch.ID
	err = s.storage.AddTxCommitment(&commitmentInBatch)
	s.NoError(err)

	return batch
}

func (s *TransferTestSuite) batchAddTransfers(transfers []models.TransferWithBatchDetails) {
	txs := make([]models.Transfer, 5)
	for i := range transfers {
		txs[i] = transfers[i].Transfer
	}
	err := s.storage.BatchAddTransfer(txs)
	s.NoError(err)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
