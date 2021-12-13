package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RegisteredSpokeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *RegisteredSpokeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RegisteredSpokeTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *RegisteredSpokeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *RegisteredSpokeTestSuite) TestAddRegisteredSpoke_AddAndRetrieve() {
	registeredSpoke := &models.RegisteredSpoke{
		ID:       models.MakeUint256(1),
		Contract: common.BytesToAddress(utils.NewRandomHash().Bytes()),
	}
	err := s.storage.AddRegisteredSpoke(registeredSpoke)
	s.NoError(err)

	actual, err := s.storage.GetRegisteredSpoke(registeredSpoke.ID)
	s.NoError(err)

	s.Equal(registeredSpoke, actual)
}

func (s *RegisteredSpokeTestSuite) TestGetRegisteredSpoke_NonexistentSpoke() {
	res, err := s.storage.GetRegisteredSpoke(models.MakeUint256(42))
	s.ErrorIs(err, NewNotFoundError("registered spoke"))
	s.Nil(res)
}

func TestRegisteredSpokeTestSuite(t *testing.T) {
	suite.Run(t, new(RegisteredSpokeTestSuite))
}
