package storage

import (
	db2 "github.com/Worldcoin/hubble-commander/db"
	"math/big"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	. "github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db2.TestDB
}

func (s *StorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StorageTestSuite) SetupTest() {
	testDB, err := db2.GetTestDB()
	s.NoError(err)
	s.storage = &Storage{DB: testDB.DB}
	s.db = testDB
}

func (s *StorageTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StorageTestSuite) TestAddTransaction() {
	tx := &models.Transaction{
		Hash:      common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		FromIndex: big.NewInt(1),
		ToIndex:   big.NewInt(2),
		Amount:    big.NewInt(1000),
		Fee:       big.NewInt(100),
		Nonce:     big.NewInt(0),
		Signature: []byte{1, 2, 3, 4, 5},
	}
	err := s.storage.AddTransaction(tx)
	s.NoError(err)

	var hash string
	err = sq.Select("*").From("transaction").
		RunWith(s.storage.DB).
		Scan(
			&hash,
			String(""),
			String(""),
			String(""),
			String(""),
			String(""),
			String(""),
		)
	s.NoError(err)
	s.Equal("0x0000000000000000000000000000000000000000000000000000000102030405", hash)

	var count int
	err = sq.Select("count(*)").From("transaction").
		RunWith(s.storage.DB).Scan(&count)
	s.NoError(err)
	s.Equal(1, count)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
