package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
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
	err := s.storage.AddTransaction(&transfer)
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
	err := s.storage.AddTransaction(&includedTransfer)
	s.NoError(err)

	res, err := s.storage.GetTransfer(transfer.Hash)
	s.NoError(err)
	s.Equal(includedTransfer, *res)
}

func (s *TransferTestSuite) TestGetTransfer_DifferentTxType() {
	err := s.storage.AddTransaction(&create2Transfer)
	s.NoError(err)

	_, err = s.storage.GetTransfer(create2Transfer.Hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
}

func (s *TransferTestSuite) TestMarkTransfersAsIncluded() {
	txs := make([]models.Transfer, 2)
	for i := 0; i < len(txs); i++ {
		txs[i] = transfer
		txs[i].Hash = utils.RandomHash()
		err := s.storage.AddTransaction(&txs[i])
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

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
