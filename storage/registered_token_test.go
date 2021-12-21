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

func (s *RegisteredTokenTestSuite) TestGetRegisteredToken_NonexistentToken() {
	res, err := s.storage.GetRegisteredToken(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("registered token"))
	s.Nil(res)
}

func TestRegisteredTokenTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredTokenTestSuite))
}
