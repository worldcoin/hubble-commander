package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.TestStorage
	tree       *st.StateTree
	commitment models.Commitment
	batch      models.Batch
}

func (s *GetBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{nil, s.storage.Storage, nil}
	s.tree = st.NewStateTree(s.storage.Storage)

	hash := utils.RandomHash()
	s.commitment = commitment
	s.commitment.IncludedInBatch = &hash
	s.commitment.AccountTreeRoot = &hash

	s.batch = models.Batch{
		Hash:              hash,
		Type:              txtype.Transfer,
		FinalisationBlock: 42000,
	}
}

func (s *GetBatchTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchTestSuite) TestGetBatchByHash() {
	err := addLeaf(s.storage, s.tree)
	s.NoError(err)
	err = s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch, result.Batch)
	s.Equal(getSubmissionBlock(s.batch.FinalisationBlock), result.SubmissionBlock)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NoCommitments() {
	err := addLeaf(s.storage, s.tree)
	s.NoError(err)
	err = s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(s.batch.Hash)
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonExistentBatch() {
	result, err := s.api.GetBatchByHash(*s.commitment.IncludedInBatch)
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID() {
	err := addLeaf(s.storage, s.tree)
	s.NoError(err)
	err = s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch, result.Batch)
	s.Equal(getSubmissionBlock(s.batch.FinalisationBlock), result.SubmissionBlock)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NoCommitments() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NonExistentBatch() {
	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func addLeaf(storage *st.Storage, tree *st.StateTree) error {
	err := storage.AddAccountIfNotExists(&models.Account{PubKeyID: 1})
	if err != nil {
		return err
	}

	return tree.Set(uint32(1), &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	})
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
