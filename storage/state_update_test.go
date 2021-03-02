package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StateUpdateTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *StateUpdateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StateUpdateTestSuite) SetupTest() {
	testDB, err := db.GetTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *StateUpdateTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *StateUpdateTestSuite) Test_AddStateUpdate_AddAndRetrieve() {
	path, err := models.NewMerklePath("00001111111111001111111111111111")
	s.NoError(err)
	update := &models.StateUpdate{
		ID:          1,
		MerklePath:  *path,
		CurrentHash: common.BytesToHash([]byte{1, 2}),
		CurrentRoot: common.BytesToHash([]byte{1, 2, 3}),
		PrevHash:    common.BytesToHash([]byte{1, 2, 3, 4}),
		PrevRoot:    common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}
	err = s.storage.AddStateUpdate(update)
	s.NoError(err)

	res, err := s.storage.GetStateUpdate(1)
	s.NoError(err)

	s.Equal(update, res)
}

func (s *StateUpdateTestSuite) Test_GetStateUpdate_NonExistentUpdate() {
	res, err := s.storage.GetStateUpdate(1)
	s.NoError(err)
	s.Nil(res)
}

func TestStateUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(StateUpdateTestSuite))
}
