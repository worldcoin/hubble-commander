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
	storage := Storage{}
	currentBlockNumber := uint32(420)

	storage.SetLatestBlockNumber(currentBlockNumber)

	latestBlockNumber := storage.GetLatestBlockNumber()

	s.Equal(currentBlockNumber, latestBlockNumber)
	s.Equal(currentBlockNumber, storage.latestBlockNumber)
}

func (s *BlockNumberTestSuite) TestGetSyncedBlock() {
	err := s.storage.SetChainState(&chainState)
	s.NoError(err)

	latestBlockNumber, err := s.storage.GetSyncedBlock(chainState.ChainID)
	s.NoError(err)

	s.Equal(chainState.SyncedBlock, *latestBlockNumber)
}

func (s *BlockNumberTestSuite) TestGetSyncedBlock_NoExistentChainState() {
	latestBlockNumber, err := s.storage.GetSyncedBlock(chainState.ChainID)
	s.NoError(err)

	s.Equal(uint32(0), *latestBlockNumber)
}

func (s *BlockNumberTestSuite) TestSetSyncedBlock() {
	err := s.storage.SetChainState(&chainState)
	s.NoError(err)

	blockNumber := uint32(450)
	err = s.storage.SetSyncedBlock(chainState.ChainID, blockNumber)
	s.NoError(err)

	latestBlockNumber, err := s.storage.GetSyncedBlock(chainState.ChainID)
	s.NoError(err)

	s.Equal(blockNumber, *latestBlockNumber)
	s.Equal(blockNumber, *s.storage.syncedBlock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(BlockNumberTestSuite))
}
