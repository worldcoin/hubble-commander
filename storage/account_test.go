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

func (s *AccountTestSuite) TestGetUnusedPubKeyID_NoAccounts() {
	_, err := s.storage.GetUnusedPubKeyID(&account1.PublicKey, models.NewUint256(100))
	s.ErrorIs(err, NewNotFoundError("account leaves"))
}

func (s *AccountTestSuite) TestGetUnusedPubKeyID_ExistingAccountAndNoUserState() {
	err := s.storage.AccountTree.SetSingle(&account1)
	s.NoError(err)
	pubKeyID, err := s.storage.GetUnusedPubKeyID(&account1.PublicKey, models.NewUint256(100))
	s.NoError(err)
	s.EqualValues(1, *pubKeyID)
}

func (s *AccountTestSuite) TestGetUnusedPubKeyID_NoUnusedPubKeyIDs() {
	account := models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	leaf := &models.StateLeaf{
		UserState: models.UserState{
			PubKeyID: 0,
			TokenID:  models.MakeUint256(1),
		},
	}
	_, err = s.storage.StateTree.Set(leaf.StateID, &leaf.UserState)
	s.NoError(err)

	_, err = s.storage.GetUnusedPubKeyID(&models.PublicKey{1, 2, 3}, &leaf.TokenID)
	s.ErrorIs(err, NewNotFoundError("pub key id"))
}

func (s *AccountTestSuite) TestGetUnusedPubKeyID_ReturnsFirstUnusedPubKeyID() {
	accounts := []models.AccountLeaf{
		{PubKeyID: 0, PublicKey: models.PublicKey{1, 2, 3}},
		{PubKeyID: 1, PublicKey: models.PublicKey{2, 3, 4}},
		{PubKeyID: 2, PublicKey: models.PublicKey{2, 3, 4}},
		{PubKeyID: 3, PublicKey: models.PublicKey{2, 3, 4}},
		{PubKeyID: 4, PublicKey: models.PublicKey{2, 3, 4}},
	}

	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	leaves := []models.StateLeaf{
		{
			StateID: 0,
			UserState: models.UserState{
				PubKeyID: 1,
				TokenID:  models.MakeUint256(1),
			},
		},
		{
			StateID: 1,
			UserState: models.UserState{
				PubKeyID: 2,
				TokenID:  models.MakeUint256(1),
			},
		},
	}

	for i := range leaves {
		_, err := s.storage.StateTree.Set(leaves[i].StateID, &leaves[i].UserState)
		s.NoError(err)
	}

	pubKeyID, err := s.storage.GetUnusedPubKeyID(&models.PublicKey{2, 3, 4}, models.NewUint256(1))
	s.NoError(err)
	s.Equal(uint32(3), *pubKeyID)
}

func (s *AccountTestSuite) TestGetUnusedPubKeyID_MultipleTokenIDs() {
	accounts := []models.AccountLeaf{
		{PubKeyID: 1, PublicKey: models.PublicKey{2, 3, 4}},
		{PubKeyID: 2, PublicKey: models.PublicKey{2, 3, 4}},
	}

	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	leaves := []models.StateLeaf{
		{
			StateID: 0,
			UserState: models.UserState{
				PubKeyID: 1,
				TokenID:  models.MakeUint256(1),
			},
		},
		{
			StateID: 1,
			UserState: models.UserState{
				PubKeyID: 2,
				TokenID:  models.MakeUint256(2),
			},
		},
	}
	for i := range leaves {
		_, err := s.storage.StateTree.Set(leaves[i].StateID, &leaves[i].UserState)
		s.NoError(err)
	}

	pubKeyID, err := s.storage.GetUnusedPubKeyID(&accounts[1].PublicKey, models.NewUint256(1))
	s.NoError(err)
	s.EqualValues(2, *pubKeyID)
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

func (s *AccountTestSuite) TestGetPublicKeyByStateID_NonExistentStateLeaf() {
	_, err := s.storage.GetPublicKeyByStateID(1)
	s.ErrorIs(err, NewNotFoundError("state leaf"))
}

func (s *AccountTestSuite) TestGetPublicKeyByStateID_NonExistentAccountLeaf() {
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
