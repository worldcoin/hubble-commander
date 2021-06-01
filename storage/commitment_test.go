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
	batch := models.Batch{
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.NewUint256(123),
		FinalisationBlock: ref.Uint32(1234),
	}
	id, err := s.storage.AddBatch(&batch)
	s.NoError(err)
	return id
}

func (s *CommitmentTestSuite) TestMarkCommitmentAsIncluded_UpdatesRecord() {
	batchID := s.addRandomBatch()

	id, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(*id, *batchID)
	s.NoError(err)

	expected := s.getCommitment(*id)
	expected.IncludedInBatch = batchID

	actual, err := s.storage.GetCommitment(*id)
	s.NoError(err)

	s.Equal(expected, actual)
}

func (s *CommitmentTestSuite) TestUpdateCommitmentAccountTreeRoot_UpdatesRecord() {
	batchID := s.addRandomBatch()
	accountRoot := utils.RandomHash()

	id1, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)
	id2, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)
	id3, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	err = s.storage.MarkCommitmentAsIncluded(*id1, *batchID)
	s.NoError(err)
	err = s.storage.MarkCommitmentAsIncluded(*id3, *batchID)
	s.NoError(err)

	err = s.storage.UpdateCommitmentsAccountTreeRoot(*batchID, accountRoot)
	s.NoError(err)

	expectedCommitment1 := s.getCommitment(*id1)
	expectedCommitment1.IncludedInBatch = batchID
	expectedCommitment1.AccountTreeRoot = &accountRoot

	expectedCommitment2 := s.getCommitment(*id2)

	expectedCommitment3 := s.getCommitment(*id3)
	expectedCommitment3.IncludedInBatch = batchID
	expectedCommitment3.AccountTreeRoot = &accountRoot

	actualCommitment1, err := s.storage.GetCommitment(*id1)
	s.NoError(err)
	actualCommitment2, err := s.storage.GetCommitment(*id2)
	s.NoError(err)
	actualCommitment3, err := s.storage.GetCommitment(*id3)
	s.NoError(err)

	s.Equal(expectedCommitment1, actualCommitment1)
	s.Equal(expectedCommitment2, actualCommitment2)
	s.Equal(expectedCommitment3, actualCommitment3)
}

func (s *CommitmentTestSuite) TestGetCommitment_NonExistentCommitment() {
	res, err := s.storage.GetCommitment(42)
	s.Equal(NewNotFoundError("commitment"), err)
	s.Nil(res)
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID() {
	_, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	batchID := s.addRandomBatch()
	includedCommitment := commitment
	includedCommitment.IncludedInBatch = batchID
	includedCommitment.FeeReceiver = 0
	includedCommitment.AccountTreeRoot = utils.NewRandomHash()

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

	s.addLeaf()

	commitments, err := s.storage.GetCommitmentsByBatchID(*batchID)
	s.NoError(err)
	s.Len(commitments, 2)
	s.Contains(commitments, expectedCommitments[0])
	s.Contains(commitments, expectedCommitments[1])
}

func (s *CommitmentTestSuite) TestGetCommitmentsByBatchID_NonExistentCommitments() {
	batchID := s.addRandomBatch()
	commitments, err := s.storage.GetCommitmentsByBatchID(*batchID)
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
