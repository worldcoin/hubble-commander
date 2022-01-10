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
	mmCommitment = models.MMCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      models.MakeUint256(1),
				IndexInBatch: 0,
			},
			Type:          batchtype.MassMigration,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       uint32(1),
		CombinedSignature: models.MakeRandomSignature(),
		BodyHash:          utils.NewRandomHash(),
		Meta: &models.MassMigrationMeta{
			SpokeID:     1,
			TokenID:     models.MakeUint256(2),
			Amount:      models.MakeUint256(3),
			FeeReceiver: 4,
		},
		WithdrawRoot: utils.RandomHash(),
	}
)

type MMCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *TestStorage
}

func (s *MMCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MMCommitmentTestSuite) SetupTest() {
	var err error
	s.storage, err = NewTestStorage()
	s.NoError(err)
}

func (s *MMCommitmentTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MMCommitmentTestSuite) TestAddMMCommitment_AddAndRetrieve() {
	err := s.storage.AddMMCommitment(&mmCommitment)
	s.NoError(err)

	actual, err := s.storage.GetMMCommitment(&mmCommitment.ID)
	s.NoError(err)
	s.Equal(mmCommitment, *actual)
}

func (s *MMCommitmentTestSuite) addRandomBatch() models.Uint256 {
	batch := models.Batch{
		ID:                models.MakeUint256(123),
		Type:              batchtype.MassMigration,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(1234),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return batch.ID
}

func (s *MMCommitmentTestSuite) TestGetMMCommitment_NonexistentCommitment() {
	res, err := s.storage.GetMMCommitment(&mmCommitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func (s *MMCommitmentTestSuite) TestGetMMCommitment_InvalidCommitmentType() {
	err := s.storage.AddTxCommitment(&txCommitment)
	s.NoError(err)

	res, err := s.storage.GetMMCommitment(&txCommitment.ID)
	s.ErrorIs(err, NewNotFoundError("commitment"))
	s.Nil(res)
}

func (s *MMCommitmentTestSuite) TestUpdateMMCommitments() {
	expectedCommitments := make([]models.MMCommitment, 2)
	for i := range expectedCommitments {
		expectedCommitments[i] = mmCommitment
		expectedCommitments[i].ID.IndexInBatch = uint8(i)

		err := s.storage.AddMMCommitment(&expectedCommitments[i])
		s.NoError(err)

		expectedCommitments[i].BodyHash = utils.NewRandomHash()
	}

	err := s.storage.UpdateMMCommitments(expectedCommitments)
	s.NoError(err)

	commitments, err := s.storage.GetMMCommitmentsByBatchID(expectedCommitments[0].ID.BatchID)
	s.NoError(err)
	s.Equal(expectedCommitments, commitments)
}

func (s *MMCommitmentTestSuite) TestUpdateMMCommitments_NonexistentCommitment() {
	commitment := mmCommitment
	commitment.BodyHash = utils.NewRandomHash()
	err := s.storage.UpdateMMCommitments([]models.MMCommitment{commitment})
	s.ErrorIs(err, NewNotFoundError("commitment"))
}

func (s *MMCommitmentTestSuite) TestGetMMCommitmentsByBatchID() {
	err := s.storage.AddMMCommitment(&mmCommitment)
	s.NoError(err)

	batchID := s.addRandomBatch()
	includedCommitment := mmCommitment
	includedCommitment.ID.BatchID = batchID

	expectedCommitments := make([]models.MMCommitment, 2)
	for i := 0; i < 2; i++ {
		includedCommitment.ID.IndexInBatch = uint8(i)
		err = s.storage.AddMMCommitment(&includedCommitment)
		s.NoError(err)

		expectedCommitments[i] = includedCommitment
	}

	commitments, err := s.storage.GetMMCommitmentsByBatchID(batchID)
	s.NoError(err)
	s.Len(commitments, 2)
	s.Contains(commitments, expectedCommitments[0])
	s.Contains(commitments, expectedCommitments[1])
}

func (s *MMCommitmentTestSuite) TestGetMMCommitmentsByBatchID_NonexistentCommitments() {
	batchID := s.addRandomBatch()
	commitments, err := s.storage.GetMMCommitmentsByBatchID(batchID)
	s.NoError(err)
	s.Len(commitments, 0)
}

func TestMMCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(MMCommitmentTestSuite))
}
