package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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

func (s *AccountTestSuite) Test_AddAccountIfNotExists_AddAndRetrieve() {
	account := models.Account{
		PubkeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account}, res)
}

func (s *AccountTestSuite) Test_GetAccounts_ReturnsAllAccounts() {
	pubKey := models.PublicKey{1, 2, 3}
	accounts := []models.Account{{
		PubkeyID:  0,
		PublicKey: pubKey,
	}, {
		PubkeyID:  1,
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

func (s *AccountTestSuite) Test_AddAccountIfNotExists_Idempotent() {
	account := models.Account{
		PubkeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account}, res)
}

func (s *AccountTestSuite) Test_GetPublicKey_ReturnsPublicKey() {
	account := models.Account{
		PubkeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	key, err := s.storage.GetPublicKey(0)
	s.NoError(err)
	s.Equal(account.PublicKey, *key)
}

func (s *AccountTestSuite) Test_GetUnusedPubKeyID() {
	accounts := []models.Account{
		{
			PubkeyID:  0,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubkeyID:  1,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubkeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
	}

	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubkeyID:   0,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err := s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	pubKeyID, err := s.storage.GetUnusedPubKeyID(&accounts[1].PublicKey)
	s.NoError(err)
	s.Equal(uint32(1), *pubKeyID)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
