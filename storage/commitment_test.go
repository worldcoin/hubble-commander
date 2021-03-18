package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		LeafHash:          utils.RandomHash(),
		PostStateRoot:     utils.RandomHash(),
		BodyHash:          utils.RandomHash(),
		AccountTreeRoot:   utils.RandomHash(),
		CombinedSignature: models.Signature{models.MakeUint256(1), models.MakeUint256(2)},
		FeeReceiver:       uint32(1),
		Transactions:      []byte{1, 2, 3},
		IncludedInBatch:   nil,
	}
)

type CommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *Storage
	db      *db.TestDB
}

func (s *CommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.storage = NewTestStorage(testDB.DB)
	s.db = testDB
}

func (s *CommitmentTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CommitmentTestSuite) Test_AddCommitment_AddAndRetrieve() {
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(commitment.LeafHash)
	s.NoError(err)
	s.Equal(commitment, *actual)
}

func (s *CommitmentTestSuite) Test_MarkCommitmentAsIncluded_UpdatesRecord() {
	batch := models.Batch{Hash: utils.RandomHash()}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	err = s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(commitment.LeafHash, batch.Hash)
	s.NoError(err)

	expected := commitment
	expected.IncludedInBatch = &batch.Hash

	actual, err := s.storage.GetCommitment(commitment.LeafHash)
	s.NoError(err)

	s.Equal(expected, *actual)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
