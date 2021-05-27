package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.TestStorage
	testClient *eth.TestClient
	tree       *st.StateTree
	commitment models.Commitment
	batch      models.Batch
}

func (s *GetBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage, client: s.testClient.Client}
	s.tree = st.NewStateTree(s.storage.Storage)

	hash := utils.RandomHash()
	s.commitment = commitment
	s.commitment.IncludedInBatch = ref.Int32(1)
	s.commitment.AccountTreeRoot = &hash

	s.batch = models.Batch{
		ID:                1,
		Hash:              &hash,
		Type:              txtype.Transfer,
		FinalisationBlock: ref.Uint32(42000),
		Number:            models.NewUint256(1),
	}
}

func (s *GetBatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchTestSuite) TestGetBatchByHash() {
	s.addLeaf()
	_, err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch, result.Batch)
	s.Equal(*s.batch.FinalisationBlock-rollup.DefaultBlocksToFinalise, result.SubmissionBlock)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NoCommitments() {
	s.addLeaf()
	_, err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonExistentBatch() {
	result, err := s.api.GetBatchByHash(utils.RandomHash())
	s.Equal(st.NewNotFoundError("batch"), err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID() {
	s.addLeaf()
	_, err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	_, err = s.storage.AddCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch, result.Batch)
	s.Equal(*s.batch.FinalisationBlock-rollup.DefaultBlocksToFinalise, result.SubmissionBlock)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NoCommitments() {
	_, err := s.storage.AddBatch(&s.batch)
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

func (s *GetBatchTestSuite) addLeaf() {
	err := s.storage.AddAccountIfNotExists(&models.Account{PubKeyID: 1})
	s.NoError(err)

	err = s.tree.Set(uint32(1), &models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
