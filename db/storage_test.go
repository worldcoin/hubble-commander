package db

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	*require.Assertions
	suite.Suite
	storage Storage
	db      *sqlx.DB
}

func (s *StorageTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	cfg := config.GetTestConfig()
	migrator, err := GetMigrator(&cfg)
	s.NoError(err)

	s.NoError(migrator.Up())
}

func (s *StorageTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	db, err := GetTestDB(&cfg)
	s.NoError(err)
	s.db = db
	s.storage = Storage{db}
}

func (s *StorageTestSuite) TearDownTest() {
	cfg := config.GetTestConfig()
	migrator, err := GetMigrator(&cfg)
	s.NoError(err)

	s.NoError(migrator.Down())

	err = s.db.Close()
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
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
