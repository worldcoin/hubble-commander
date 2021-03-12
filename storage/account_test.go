package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func (s *AccountTestSuite) Test_AddAccount_AddAndRetrieve() {
	account := models.Account{
		AccountIndex: 0,
		PublicKey:    [128]byte{},
	}

	err := s.storage.AddAccount(&account)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account}, res)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
