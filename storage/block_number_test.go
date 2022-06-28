package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BlockNumberTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *BlockNumberTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *BlockNumberTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *BlockNumberTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *BlockNumberTestSuite) TestSetLatestBlockNumber() {
	currentBlockNumber := uint32(420)

	s.storage.SetLatestBlockNumber(currentBlockNumber)

	latestBlockNumber := s.storage.GetLatestBlockNumber()

	s.Equal(currentBlockNumber, latestBlockNumber)
	s.Equal(currentBlockNumber, s.storage.latestBlockNumber)
}

func (s *BlockNumberTestSuite) TestGetSyncedBlock() {
	err := s.storage.SetSyncedBlock(10)
	s.NoError(err)

	syncedBlock, err := s.storage.GetSyncedBlock()
	s.NoError(err)

	s.Equal(uint64(10), *syncedBlock)
}

func (s *BlockNumberTestSuite) TestGetSyncedBlock_NoExistentChainState() {
	syncedBlock, err := s.storage.GetSyncedBlock()
	s.NoError(err)

	s.Equal(uint64(0), *syncedBlock)
}

func (s *BlockNumberTestSuite) TestSetSyncedBlock() {
	blockNumber := uint64(450)
	err := s.storage.SetSyncedBlock(blockNumber)
	s.NoError(err)

	syncedBlock, err := s.storage.GetSyncedBlock()
	s.NoError(err)

	s.Equal(blockNumber, *syncedBlock)
	s.Equal(blockNumber, s.storage.syncedBlock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(BlockNumberTestSuite))
}
