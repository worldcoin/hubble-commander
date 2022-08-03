package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
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
			CommitmentSlot: &models.CommitmentSlot{
				BatchID: models.MakeUint256(1),
				IndexInBatch: 0,
				IndexInCommitment: 0,
			},
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

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_DifferentTxType() {
	err := s.storage.AddTransaction(&transfer)
	s.NoError(err)

	_, err = s.storage.GetCreate2Transfer(transfer.Hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
}

func (s *Create2TransferTestSuite) TestGetCreate2Transfer_NonexistentTransaction() {
	hash := common.BytesToHash([]byte{1, 2, 3, 4, 5})
	res, err := s.storage.GetCreate2Transfer(hash)
	s.ErrorIs(err, NewNotFoundError("transaction"))
	s.Nil(res)
}

func TestCreate2TransferTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferTestSuite))
}
