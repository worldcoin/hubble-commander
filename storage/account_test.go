package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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
	storage *TestStorage
}

func (s *AccountTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
}

func (s *AccountTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *AccountTestSuite) TestAddAccountIfNotExists_AddAndRetrieve() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	res, err := s.storage.GetAccounts(&account1.PublicKey)
	s.NoError(err)

	s.Equal([]models.Account{account1}, res)
}

func (s *AccountTestSuite) TestGetAccounts_NoPublicKeys() {
	_, err := s.storage.GetAccounts(&models.PublicKey{1, 2, 3})
	s.Equal(NewNotFoundError("accounts"), err)
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

func (s *AccountTestSuite) Test_GetUnusedPubKeyID_NoPublicKeys() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)
	_, err = s.storage.GetUnusedPubKeyID(&account1.PublicKey)
	s.Equal(NewNotFoundError("pub key id"), err)
}

func (s *AccountTestSuite) Test_GetUnusedPubKeyID_NoUnusedPublicIDs() {
	account := models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   0,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err = s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	_, err = s.storage.GetUnusedPubKeyID(&models.PublicKey{1, 2, 3})
	s.Equal(NewNotFoundError("pub key id"), err)
}

func (s *AccountTestSuite) Test_GetUnusedPubKeyID() {
	accounts := []models.Account{
		{
			PubKeyID:  0,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  3,
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
			PubKeyID:   1,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	leaf2 := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
		UserState: models.UserState{
			PubKeyID:   2,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err := s.storage.AddStateLeaf(leaf)
	s.NoError(err)
	err = s.storage.AddStateLeaf(leaf2)
	s.NoError(err)

	pubKeyID, err := s.storage.GetUnusedPubKeyID(&accounts[1].PublicKey)
	s.NoError(err)
	s.Equal(uint32(3), *pubKeyID)
}

func (s *AccountTestSuite) Test_GetUnusedPubKeyID_MultipleTokenIndexes() {
	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{2, 3, 4},
		},
	}

	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	leaves := []models.StateLeaf{
		{
			DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
			UserState: models.UserState{
				PubKeyID:   1,
				TokenIndex: models.MakeUint256(1),
				Balance:    models.MakeUint256(420),
				Nonce:      models.MakeUint256(0),
			},
		},
		{
			DataHash: common.BytesToHash([]byte{2, 3, 4, 5, 6}),
			UserState: models.UserState{
				PubKeyID:   2,
				TokenIndex: models.MakeUint256(2),
				Balance:    models.MakeUint256(420),
				Nonce:      models.MakeUint256(0),
			},
		},
	}
	for i := range leaves {
		err := s.storage.AddStateLeaf(&leaves[i])
		s.NoError(err)
	}

	pubKeyID, err := s.storage.GetUnusedPubKeyID(&accounts[1].PublicKey, leaves[1].TokenIndex)
	s.Equal(NewNotFoundError("pub key id"), err)
	s.Nil(pubKeyID)
}

func (s *AccountTestSuite) Test_DoesAccountExist_Exists() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	exists, err := s.storage.AccountExists(&account1.PublicKey)
	s.NoError(err)
	s.True(exists)
}

func (s *AccountTestSuite) TestAccountExists_NotExists() {
	exists, err := s.storage.AccountExists(&account1.PublicKey)
	s.NoError(err)
	s.False(exists)
}

func (s *AccountTestSuite) TestAccountWithTokenExists_Exists() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	tokenIndex := models.MakeUint256(5)
	leaf := &models.StateLeaf{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
		UserState: models.UserState{
			PubKeyID:   1,
			TokenIndex: tokenIndex,
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
	}
	err = s.storage.AddStateLeaf(leaf)
	s.NoError(err)

	exists, err := s.storage.AccountWithTokenExists(&account1.PublicKey, tokenIndex)
	s.NoError(err)
	s.True(exists)
}

func (s *AccountTestSuite) TestAccountWithTokenExists_AccountNotExists() {
	exists, err := s.storage.AccountWithTokenExists(&account1.PublicKey, models.MakeUint256(1))
	s.NoError(err)
	s.False(exists)
}

func (s *AccountTestSuite) TestAccountWithTokenExists_StateLeafWithTokenNotExists() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	exists, err := s.storage.AccountWithTokenExists(&account1.PublicKey, models.MakeUint256(1))
	s.NoError(err)
	s.False(exists)
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
