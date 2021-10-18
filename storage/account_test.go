package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	account1 = models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{3, 4, 5},
	}
	account2 = models.AccountLeaf{
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
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *AccountTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *AccountTestSuite) TestGetFirstPubKeyID_NoAccounts() {
	_, err := s.storage.GetFirstPubKeyID(&account1.PublicKey)
	s.ErrorIs(err, NewNotFoundError("pub key id"))
}

func (s *AccountTestSuite) TestGetFirstPubKeyID_ExistingAccount() {
	err := s.storage.AccountTree.SetSingle(&account1)
	s.NoError(err)
	pubKeyID, err := s.storage.GetFirstPubKeyID(&account1.PublicKey)
	s.NoError(err)
	s.EqualValues(account1.PubKeyID, *pubKeyID)
}

func (s *AccountTestSuite) TestGetFirstPubKeyID_MultipleAccounts() {
	accounts := []models.AccountLeaf{
		{PubKeyID: 1, PublicKey: models.PublicKey{2, 3, 4}},
		{PubKeyID: 2, PublicKey: models.PublicKey{2, 3, 4}},
	}

	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	pubKeyID, err := s.storage.GetFirstPubKeyID(&accounts[0].PublicKey)
	s.NoError(err)
	s.EqualValues(accounts[0].PubKeyID, *pubKeyID)
}

func (s *AccountTestSuite) TestGetPublicKeyByStateID() {
	err := s.storage.AccountTree.SetSingle(&account1)
	s.NoError(err)
	err = s.storage.AccountTree.SetSingle(&account2)
	s.NoError(err)

	leaves := []models.StateLeaf{
		{
			StateID: 1,
			UserState: models.UserState{
				PubKeyID: 1,
				TokenID:  models.MakeUint256(1),
			},
		},
		{
			StateID: 2,
			UserState: models.UserState{
				PubKeyID: 2,
				TokenID:  models.MakeUint256(2),
			},
		},
	}
	for i := range leaves {
		_, setErr := s.storage.StateTree.Set(leaves[i].StateID, &leaves[i].UserState)
		s.NoError(setErr)
	}

	publicKey, err := s.storage.GetPublicKeyByStateID(2)
	s.NoError(err)
	s.Equal(account2.PublicKey, *publicKey)
}

func (s *AccountTestSuite) TestGetPublicKeyByStateID_NonexistentStateLeaf() {
	_, err := s.storage.GetPublicKeyByStateID(1)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
}

func (s *AccountTestSuite) TestGetPublicKeyByStateID_NonexistentAccountLeaf() {
	userState := &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
	}

	_, err := s.storage.StateTree.Set(1, userState)
	s.NoError(err)

	_, err = s.storage.GetPublicKeyByStateID(1)
	s.ErrorIs(err, NewNotFoundError("account leaf"))
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}
