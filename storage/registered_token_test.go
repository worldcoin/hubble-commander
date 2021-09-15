package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisteredTokenTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *RegisteredTokenTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisteredTokenTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *RegisteredTokenTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *RegisteredTokenTestSuite) TestAddRegisteredToken_AddAndRetrieve() {
	registeredToken := &models.RegisteredToken{
		ID:       models.MakeUint256(1),
		Contract: common.BytesToAddress(utils.NewRandomHash().Bytes()),
	}
	err := s.storage.AddRegisteredToken(registeredToken)
	s.NoError(err)

	actual, err := s.storage.GetRegisteredToken(registeredToken.ID)
	s.NoError(err)

	s.Equal(registeredToken, actual)
}

func (s *RegisteredTokenTestSuite) TestGetRegisteredToken_NonExistentToken() {
	res, err := s.storage.GetRegisteredToken(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("registered token"))
	s.Nil(res)
}

func (s *RegisteredTokenTestSuite) TestDeleteRegisteredTokens() {
	registeredTokens := []models.RegisteredToken{
		{
			ID:       models.MakeUint256(1),
			Contract: common.BytesToAddress(utils.NewRandomHash().Bytes()),
		},
		{
			ID:       models.MakeUint256(2),
			Contract: common.BytesToAddress(utils.NewRandomHash().Bytes()),
		},
	}
	for i := range registeredTokens {
		err := s.storage.AddRegisteredToken(&registeredTokens[i])
		s.NoError(err)
	}

	err := s.storage.DeleteRegisteredTokens(registeredTokens[0].ID, registeredTokens[1].ID)
	s.NoError(err)

	for i := range registeredTokens {
		_, err := s.storage.GetRegisteredToken(registeredTokens[i].ID)
		s.ErrorIs(err, NewNotFoundError("registered token"))
	}
}

func (s *RegisteredTokenTestSuite) TestDeleteRegisteredTokens_NonExistentToken() {
	err := s.storage.DeleteRegisteredTokens(models.MakeUint256(1))
	s.ErrorIs(err, NewNotFoundError("registered token"))
}

func TestRegisteredTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredTokenTestSuite))
}
