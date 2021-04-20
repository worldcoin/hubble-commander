package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	account1 = models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{3, 4, 5},
	}
	account2 = models.Account{
		PubKeyID:  2,
		PublicKey: models.PublicKey{4, 5, 6},
	}
)

type AccountTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *AccountTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *AccountTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *AccountTestSuite) TestAddAccountIfNotExists_AddAndRetrieve() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account1.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account1}, res)
}

func (s *AccountTestSuite) TestGetAccounts_ReturnsAllAccounts() {
	pubKey := models.PublicKey{1, 2, 3}
	accounts := []models.Account{{
		PubKeyID:  0,
		PublicKey: pubKey,
	}, {
		PubKeyID:  1,
		PublicKey: pubKey,
	}}

	err := s.storage.AddAccountIfNotExists(&accounts[0])
	s.NoError(err)
	err = s.storage.AddAccountIfNotExists(&accounts[1])
	s.NoError(err)

	res, err := s.storage.GetAccounts(&pubKey)
	s.NoError(err)

	s.Equal(accounts, res)
}

func (s *AccountTestSuite) TestAddAccountIfNotExists_Idempotent() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account1.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account1}, res)
}

func (s *AccountTestSuite) TestGetPublicKey_ReturnsPublicKey() {
	account := models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	key, err := s.storage.GetPublicKey(0)
	s.NoError(err)
	s.Equal(account.PublicKey, *key)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
