package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
		AccountTreeRoot:   nil,
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

func (s *CommitmentTestSuite) getCommitment(id int32) *models.Commitment {
	clone := commitment
	clone.ID = id
	return &clone
}

func (s *CommitmentTestSuite) TestAddCommitment_AddAndRetrieve() {
	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(*id)
	s.NoError(err)
	s.Equal(s.getCommitment(*id), actual)
}

func (s *CommitmentTestSuite) addRandomBatch() *common.Hash {
	batch := models.Batch{Hash: utils.RandomHash()}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return &batch.Hash
}

func (s *CommitmentTestSuite) TestMarkCommitmentAsIncluded_UpdatesRecord() {
	batchHash := s.addRandomBatch()
	accountRoot := utils.RandomHash()

	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(*id, batchHash, &accountRoot)
	s.NoError(err)

	expected := s.getCommitment(*id)
	expected.IncludedInBatch = batchHash
	expected.AccountTreeRoot = &accountRoot

	actual, err := s.storage.GetCommitment(*id)
	s.NoError(err)

	s.Equal(expected, actual)
}

func (s *CommitmentTestSuite) TestGetPendingCommitments_ReturnsOnlyPending() {
	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = s.addRandomBatch()
	_, err = s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	actual, err := s.storage.GetPendingCommitments(10)
	s.NoError(err)

	expected := commitment
	expected.ID = *id

	s.Equal([]models.Commitment{expected}, actual)
}

func (s *CommitmentTestSuite) TestGetPendingCommitments_ReturnsOnlyGivenNumberOfRows() {
	for i := 0; i < 3; i++ {
		_, err := s.storage.AddCommitment(&commitment)
		s.NoError(err)
	}

	commitments, err := s.storage.GetPendingCommitments(2)
	s.NoError(err)
	s.Len(commitments, 2)
}

func (s *StateUpdateTestSuite) TestGetCommitment_NonExistentCommitment() {
	res, err := s.storage.GetCommitment(42)
	s.Equal(NewNotFoundError("commitment"), err)
	s.Nil(res)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
