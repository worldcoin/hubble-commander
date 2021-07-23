package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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
	s.tree = NewStateTree(s.storage.StorageBase)
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

func (s *CommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func (s *CommitmentTestSuite) TestMarkCommitmentAsIncluded_UpdatesRecord() {
	batchID := s.addRandomBatch()

	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(*id, batchID)
	s.NoError(err)

	expected := s.getCommitment(*id)
	expected.IncludedInBatch = &batchID

	actual, err := s.storage.GetCommitment(*id)
	s.NoError(err)

	s.Equal(expected, actual)
}

func (s *CommitmentTestSuite) TestGetCommitment_NonExistentCommitment() {
	res, err := s.storage.GetCommitment(42)
	s.Equal(NewNotFoundError("commitment"), err)
	s.Nil(res)
}

func (s *CommitmentTestSuite) TestGetLatestCommitment() {
	expected := commitment
	for i := 0; i < 2; i++ {
		commitmentID, err := s.storage.AddCommitment(&commitment)
		s.NoError(err)
		expected.ID = *commitmentID
	}
	latestCommitment, err := s.storage.GetLatestCommitment()
	s.NoError(err)
	s.Equal(expected, *latestCommitment)
}

func (s *CommitmentTestSuite) TestGetLatestCommitment_NoCommitments() {
	_, err := s.storage.GetLatestCommitment()
	s.Equal(NewNotFoundError("commitment"), err)
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID() {
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	batchID := s.addRandomBatch()
	includedCommitment := commitment
	includedCommitment.IncludedInBatch = &batchID
	includedCommitment.FeeReceiver = 0

	expectedCommitments := make([]models.CommitmentWithTokenID, 2)
	for i := 0; i < 2; i++ {
		var commitmentID *int32
		commitmentID, err = s.storage.AddCommitment(&includedCommitment)
		s.NoError(err)
		expectedCommitments[i] = models.CommitmentWithTokenID{
			ID:                 *commitmentID,
			Transactions:       includedCommitment.Transactions,
			TokenID:            models.MakeUint256(1),
			FeeReceiverStateID: includedCommitment.FeeReceiver,
			CombinedSignature:  includedCommitment.CombinedSignature,
			PostStateRoot:      includedCommitment.PostStateRoot,
		}
	}

	s.addStateLeaf()

	commitments, err := s.storage.GetCommitmentsByBatchID(batchID)
	s.NoError(err)
	s.Len(commitments, 2)
	s.Contains(commitments, expectedCommitments[0])
	s.Contains(commitments, expectedCommitments[1])
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_NonExistentCommitments() {
	batchID := s.addRandomBatch()
	commitments, err := s.storage.GetCommitmentsByBatchID(batchID)
	s.Equal(NewNotFoundError("commitments"), err)
	s.Nil(commitments)
}

func (s *CommitmentTestSuite) TestDeleteCommitmentsByBatchIDs() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(111),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                models.MakeUint256(5),
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2345),
		},
	}
	for i := range batches {
		err := s.storage.AddBatch(&batches[i])
		s.NoError(err)

		for j := 0; j < 2; j++ {
			commitmentInBatch := commitment
			commitmentInBatch.IncludedInBatch = &batches[i].ID
			_, err = s.storage.AddCommitment(&commitmentInBatch)
			s.NoError(err)
		}
	}

	err := s.storage.DeleteCommitmentsByBatchIDs(batches[0].ID, batches[1].ID)
	s.NoError(err)
	for i := range batches {
		_, err = s.storage.GetCommitmentsByBatchID(batches[i].ID)
		s.Equal(NewNotFoundError("commitments"), err)
	}
}

func (s *CommitmentTestSuite) TestDeleteCommitmentsByBatchIDs_NoCommitments() {
	batchID := s.addRandomBatch()
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.DeleteCommitmentsByBatchIDs(batchID)
	s.Equal(ErrNoRowsAffected, err)

	_, err = s.storage.GetCommitment(*commitmentID)
	s.NoError(err)
}

func (s *CommitmentTestSuite) addStateLeaf() {
	err := s.storage.AddAccountLeafIfNotExists(&account1)
	s.NoError(err)

	_, err = s.tree.Set(uint32(0), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentTestSuite))
}
