package storage

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	transfer = models.Transfer{
		Hash:                 common.BigToHash(big.NewInt(1234)),
		FromStateID:          1,
		ToStateID:            2,
		Amount:               models.MakeUint256(1000),
		Fee:                  models.MakeUint256(100),
		Nonce:                models.MakeUint256(0),
		Signature:            []byte{1, 2, 3, 4, 5},
		IncludedInCommitment: nil,
	}
)

type TransferTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *TransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *TransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransferTestSuite) Test_AddTransfer_AddAndRetrieve() {
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	res, err := s.storage.GetTransfer(tx.Hash)
	s.NoError(err)

	s.Equal(transfer, *res)
}

func TestTransferTestSuite(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}
