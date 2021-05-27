package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      []byte{1, 2, 3},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		AccountTreeRoot:   nil,
		IncludedInBatch:   nil,
	}
)

type CommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
	tree    *StateTree
}

func (s *CommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = NewStateTree(s.storage.Storage)
}

func (s *CommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
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

func (s *CommitmentTestSuite) addRandomBatch() *int32 {
	batchHash := utils.RandomHash()
	batch := models.Batch{ID: 1, Hash: &batchHash, Number: models.NewUint256(1)}
	_, err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return &batch.ID
}

func (s *CommitmentTestSuite) TestMarkCommitmentAsIncluded_UpdatesRecord() {
	batchID := s.addRandomBatch()
	accountRoot := utils.RandomHash()

	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(*id, *batchID, &accountRoot)
	s.NoError(err)

	expected := s.getCommitment(*id)
	expected.IncludedInBatch = batchID
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

func (s *CommitmentTestSuite) TestGetCommitment_NonExistentCommitment() {
	res, err := s.storage.GetCommitment(42)
	s.Equal(NewNotFoundError("commitment"), err)
	s.Nil(res)
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchNumber() {
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	commitmentWithHash := commitment
	commitmentWithHash.FeeReceiver = 0
	commitmentWithHash.IncludedInBatch = s.addRandomBatch()
	for i := 0; i < 3; i++ {
		_, err = s.storage.AddCommitment(&commitmentWithHash)
		s.NoError(err)
	}

	s.addLeaf()

	commitments, err := s.storage.GetCommitmentsByBatchNumber(models.MakeUint256(1))
	s.NoError(err)
	s.Len(commitments, 3)
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_NonExistentCommitments() {
	_ = s.addRandomBatch()
	commitments, err := s.storage.GetCommitmentsByBatchNumber(models.MakeUint256(0))
	s.Equal(NewNotFoundError("commitments"), err)
	s.Nil(commitments)
}

func (s *CommitmentTestSuite) addLeaf() {
	err := s.storage.AddAccountIfNotExists(&account1)
	s.NoError(err)

	err = s.tree.Set(uint32(0), &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
