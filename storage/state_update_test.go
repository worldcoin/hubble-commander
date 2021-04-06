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
	testDB, err := db.NewTestDB()
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

	res, err := s.storage.GetStateUpdateByID(1)
	s.NoError(err)
	s.Equal(update, res)

	res, err = s.storage.GetStateUpdateByRoot(common.BytesToHash([]byte{1, 2, 3}))
	s.NoError(err)
	s.Equal(update, res)
}

func (s *StateUpdateTestSuite) Test_GetStateUpdate_NonExistentUpdate() {
	res, err := s.storage.GetStateUpdateByID(1)
	s.EqualError(err, "state update not found")
	s.Nil(res)
}

func (s *StateUpdateTestSuite) Test_GetLatestStateUpdate() {
	path, err := models.NewMerklePath("00001111111111001111111111111111")
	s.NoError(err)
	updates := []models.StateUpdate{
		{
			ID:          1,
			MerklePath:  *path,
			CurrentHash: common.BytesToHash([]byte{1}),
			CurrentRoot: common.BytesToHash([]byte{1}),
			PrevHash:    common.BytesToHash([]byte{1}),
			PrevRoot:    common.BytesToHash([]byte{2}),
		},
		{
			ID:          2,
			MerklePath:  *path,
			CurrentHash: common.BytesToHash([]byte{2}),
			CurrentRoot: common.BytesToHash([]byte{2}),
			PrevHash:    common.BytesToHash([]byte{2}),
			PrevRoot:    common.BytesToHash([]byte{2}),
		},
	}
	err = s.storage.AddStateUpdate(&updates[0])
	s.NoError(err)
	err = s.storage.AddStateUpdate(&updates[1])
	s.NoError(err)

	res, err := s.storage.GetLatestStateUpdate()
	s.NoError(err)
	s.Equal(&updates[1], res)
}

func TestStateUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(StateUpdateTestSuite))
}
