package storage

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		ID: models.CommitmentID{
			BatchID:      models.MakeUint256(1),
			IndexInBatch: 0,
		},
		Type:              batchtype.Transfer,
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		Transactions:      []byte{1, 2, 3},
	}
)

type CommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *CommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *CommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *CommitmentTestSuite) TestAddCommitment_AddAndRetrieve() {
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	actual, err := s.storage.GetCommitment(&commitment.ID)
	s.NoError(err)
	s.Equal(commitment, *actual)
}

func (s *CommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func (s *CommitmentTestSuite) TestGetCommitment_NonExistentCommitment() {
	res, err := s.storage.GetCommitment(&commitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func (s *CommitmentTestSuite) TestGetLatestCommitment() {
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
	s.Equal(expected, *latestCommitment)
}

func (s *CommitmentTestSuite) TestGetLatestCommitment_NoCommitments() {
	_, err := s.storage.GetLatestCommitment()
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID() {
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	batchID := s.addRandomBatch()
	includedCommitment := commitment
	includedCommitment.ID.BatchID = batchID
	includedCommitment.FeeReceiver = 0

	expectedCommitments := make([]models.CommitmentWithTokenID, 2)
	for i := 0; i < 2; i++ {
		includedCommitment.ID.IndexInBatch = uint8(i)
		err = s.storage.AddCommitment(&includedCommitment)
		s.NoError(err)

		expectedCommitments[i] = models.CommitmentWithTokenID{
			ID:                 includedCommitment.ID,
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
	s.NoError(err)
	s.Len(commitments, 0)
}

func (s *CommitmentTestSuite) TestDeleteCommitmentsByBatchIDs() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(111),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                models.MakeUint256(5),
			Type:              batchtype.Create2Transfer,
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
			commitmentInBatch.ID.BatchID = batches[i].ID
			commitmentInBatch.ID.IndexInBatch = uint8(j)
			err = s.storage.AddCommitment(&commitmentInBatch)
			s.NoError(err)
		}
	}

	err := s.storage.DeleteCommitmentsByBatchIDs(batches[0].ID, batches[1].ID)
	s.NoError(err)
	for i := range batches {
		commitments, err := s.storage.GetCommitmentsByBatchID(batches[i].ID)
		s.NoError(err)
		s.Len(commitments, 0)
	}
}

func (s *CommitmentTestSuite) TestDeleteCommitmentsByBatchIDs_NoCommitments() {
	batchID := s.addRandomBatch()
	err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.DeleteCommitmentsByBatchIDs(batchID)
	s.ErrorIs(err, NewNotFoundError("commitments"))

	_, err = s.storage.GetCommitment(&commitment.ID)
	s.NoError(err)
}

func (s *CommitmentTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(0), &models.UserState{
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
