package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TransactionTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *TransactionTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransactionTestSuite) SetupTest() {
	testDB, err := db.GetTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *TransactionTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TransactionTestSuite) TestAddTransaction() {
	tx := &models.Transaction{
		Hash:      common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		FromIndex: models.MakeUint256(1),
		ToIndex:   models.MakeUint256(2),
		Amount:    models.MakeUint256(1000),
		Fee:       models.MakeUint256(100),
		Nonce:     models.MakeUint256(0),
		Signature: []byte{1, 2, 3, 4, 5},
	}
	err := s.storage.AddTransaction(tx)
	s.NoError(err)

	res, err := s.storage.GetTransaction(tx.Hash)
	s.NoError(err)

	s.Equal(tx, res)
}

func (s *TransactionTestSuite) TestAddStateNode() {
	node := &models.StateNode{
		MerklePath: "0000111",
		DataHash:   common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err := s.storage.AddStateNode(node)
	s.NoError(err)

	res, err := s.storage.GetStateNode(node.DataHash)
	s.NoError(err)

	s.Equal(node, res)
}

func (s *TransactionTestSuite) TestAddStateLeaf() {
	leaf := &models.StateLeaf{
		DataHash:     common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		AccountIndex: models.MakeUint256(1),
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}
	err := s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	res, err := s.storage.GetStateLeaf(leaf.DataHash)
	s.NoError(err)

	s.Equal(leaf, res)
}

func (s *TransactionTestSuite) TestAddStateUpdate() {
	update := &models.StateUpdate{
		ID:          1,
		MerklePath:  "11111111111111111111111111111111",
		CurrentHash: common.BytesToHash([]byte{1, 2}),
		CurrentRoot: common.BytesToHash([]byte{1, 2, 3}),
		PrevHash:    common.BytesToHash([]byte{1, 2, 3, 4}),
		PrevRoot:    common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err := s.storage.AddStateUpdate(update)
	s.NoError(err)

	res, err := s.storage.GetStateUpdate(1)
	s.NoError(err)

	s.Equal(update, res)
}

func TestTransactionTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionTestSuite))
}
