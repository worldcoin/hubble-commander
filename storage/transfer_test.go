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
			CommitmentSlot: &models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
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
	transfer1 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 1,
			Amount:      models.MakeUint256(1000),
			Fee:         models.MakeUint256(100),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
			CommitmentSlot: &models.CommitmentSlot{
				BatchID:           models.MakeUint256(1),
				IndexInBatch:      0,
				IndexInCommitment: 0,
			},
		},
		ToStateID: 2,
	}

	err := s.storage.AddTransaction(&transfer1)
	s.NoError(err)

	expected := transfer1

	res, err := s.storage.GetTransfer(transfer1.Hash)
	s.NoError(err)
	s.Equal(expected, *res)
}

func (s *TransferTestSuite) TestAddTransfer_AddAndRetrieveIncludedTransfer() {
	includedTransfer := transfer
	includedTransfer.CommitmentSlot = &models.CommitmentSlot{
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

func (s *TransferTestSuite) TestGetTransfer_NonexistentTransfer() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetTransfer(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
