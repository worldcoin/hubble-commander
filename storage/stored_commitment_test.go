package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StoredCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *StoredCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *StoredCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment() {
	expected := commitment
	for i := 0; i < 2; i++ {
		expected.ID.IndexInBatch = uint8(i)
		err := s.storage.AddCommitment(&expected)
		s.NoError(err)
	}

	expected.ID.BatchID = models.MakeUint256(5)
	for i := 0; i < 2; i++ {
		expected.ID.IndexInBatch = uint8(i)
		err := s.storage.AddCommitment(&expected)
		s.NoError(err)
	}

	latestCommitment, err := s.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(models.MakeStoredCommitmentFromTxCommitment(&expected), *latestCommitment)
}

func (s *StoredCommitmentTestSuite) TestGetLatestCommitment_NoCommitments() {
	_, err := s.storage.GetLatestCommitment()
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func TestStoredCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(StoredCommitmentTestSuite))
}
