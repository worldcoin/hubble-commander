package storage

import (
	"math/big"
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
	create2Transfer = models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1234)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToPublicKey: account2.PublicKey,
	}
)

type Create2TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *Create2TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&account2)
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *Create2TransferTestSuite) TestAddCreate2Transfer_AddAndRetrieve() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	expected := create2Transfer

	res, err := s.storage.GetCreate2Transfer(create2Transfer.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_DifferentTxType() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	_, err = s.storage.GetCreate2Transfer(transfer.Hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
}

func (s *Create2TransferTestSuite) TestMarkCreate2TransfersAsIncluded() {
	commitmentID := &models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 1,
	}

	txs := make([]models.Create2Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = create2Transfer
		txs[i].Hash = utils.RandomHash()
		err := s.storage.AddCreate2Transfer(&txs[i])
		s.NoError(err)

		txs[i].ToStateID = ref.Uint32(uint32(i))
		txs[i].CommitmentID = commitmentID
	}

	err := s.storage.MarkCreate2TransfersAsIncluded(txs, commitmentID)
	s.NoError(err)

	for i := range txs {
		tx, err := s.storage.GetCreate2Transfer(txs[i].Hash)
		s.NoError(err)
		s.Equal(txs[i], *tx)
	}
}

func (s *Create2TransferTestSuite) TestGetCreate2TransferWithBatchDetails() {
	batch := &models.Batch{
		ID:              models.MakeUint256(1),
		Type:            batchtype.Create2Transfer,
		TransactionHash: utils.RandomHash(),
		Hash:            utils.NewRandomHash(),
		SubmissionTime:  &models.Timestamp{Time: time.Unix(170, 0).UTC()},
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitmentInBatch := txCommitment
	commitmentInBatch.ID.BatchID = batch.ID
	commitmentInBatch.ID.IndexInBatch = 0
	err = s.storage.AddTxCommitment(&commitmentInBatch)
	s.NoError(err)

	transferInBatch := create2Transfer
	transferInBatch.CommitmentID = &txCommitment.ID
	err = s.storage.AddCreate2Transfer(&transferInBatch)
	s.NoError(err)

	expected := models.Create2TransferWithBatchDetails{
		Create2Transfer: transferInBatch,
		BatchHash:       batch.Hash,
		BatchTime:       batch.SubmissionTime,
	}
	res, err := s.storage.GetCreate2TransferWithBatchDetails(transferInBatch.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransferWithBatchDetails_WithoutBatch() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	expected := models.Create2TransferWithBatchDetails{Create2Transfer: create2Transfer}

	res, err := s.storage.GetCreate2TransferWithBatchDetails(create2Transfer.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
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
	s.ErrorIs(err, ErrNoRowsAffected)
}

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_NonexistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetCreate2Transfer(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func (s *Create2TransferTestSuite) TestGetPendingCreate2Transfers() {
	commitment := &models.TxCommitment{}
	err := s.storage.AddTxCommitment(commitment)
	s.NoError(err)

	create2Transfer2 := create2Transfer
	create2Transfer2.Hash = utils.RandomHash()
	create2Transfer3 := create2Transfer
	create2Transfer3.Hash = utils.RandomHash()
	create2Transfer3.CommitmentID = &commitment.ID
	create2Transfer4 := create2Transfer
	create2Transfer4.Hash = utils.RandomHash()
	create2Transfer4.ErrorMessage = ref.String("A very boring error message")

	create2transfers := []models.Create2Transfer{
		create2Transfer,
		create2Transfer2,
		create2Transfer3,
		create2Transfer4,
	}

	err = s.storage.BatchAddCreate2Transfer(create2transfers)
	s.NoError(err)

	res, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)

	s.Len(res, 2)
	s.Contains(res, create2Transfer)
	s.Contains(res, create2Transfer2)
}

func (s *Create2TransferTestSuite) TestGetPendingCreate2Transfers_OrdersTransfersByNonceAndTxHashAscending() {
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
	create2Transfer4.Hash = common.Hash{66, 66, 66, 66}
	create2Transfer5 := create2Transfer
	create2Transfer5.TransactionBase.Nonce = models.MakeUint256(5)
	create2Transfer5.Hash = common.Hash{65, 65, 65, 65}

	create2transfers := []models.Create2Transfer{
		create2Transfer,
		create2Transfer2,
		create2Transfer3,
		create2Transfer4,
		create2Transfer5,
	}

	err := s.storage.BatchAddCreate2Transfer(create2transfers)
	s.NoError(err)

	res, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)

	s.Equal([]models.Create2Transfer{create2Transfer, create2Transfer2, create2Transfer5, create2Transfer4, create2Transfer3}, res)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByPublicKey() {
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(400),
	})
	s.NoError(err)

	transfers, err := s.storage.GetCreate2TransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(transfers, 1)
}

func (s *Create2TransferTestSuite) Test_GetCreate2TransfersByPublicKey_NoCreate2TransfersUnregisteredAccount() {
	transfers, err := s.storage.GetCreate2TransfersByPublicKey(&models.PublicKey{9, 9, 9})
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByPublicKey_NoCreate2TransfersRegisteredAccount() {
	transfers, err := s.storage.GetCreate2TransfersByPublicKey(&account2.PublicKey)
	s.NoError(err)
	s.Len(transfers, 0)
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByCommitmentID() {
	transfer1 := create2Transfer
	transfer1.CommitmentID = &txCommitment.ID
	err := s.storage.AddCreate2Transfer(&transfer1)
	s.NoError(err)

	otherCommitmentID := txCommitment.ID
	otherCommitmentID.IndexInBatch += 1
	transfer2 := create2Transfer
	transfer2.Hash = utils.RandomHash()
	transfer2.CommitmentID = &otherCommitmentID
	err = s.storage.AddCreate2Transfer(&transfer2)
	s.NoError(err)

	transfers, err := s.storage.GetCreate2TransfersByCommitmentID(&txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 1)
	s.Equal(transfer1, transfers[0])
}

func (s *Create2TransferTestSuite) TestGetCreate2TransfersByCommitmentID_NoCreate2Transfers() {
	transfers, err := s.storage.GetCreate2TransfersByCommitmentID(&txCommitment.ID)
	s.NoError(err)
	s.Len(transfers, 0)
}

func TestCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferTestSuite))
}
