package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	api                 *API
	storage             *st.TestStorage
	testClient          *eth.TestClient
	commitment          models.TxCommitment
	batch               models.Batch
	batchNotFoundAPIErr *APIError
}

func (s *GetBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage, client: s.testClient.Client}

	s.batch = models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(42000),
		AccountTreeRoot:   utils.NewRandomHash(),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}

	s.commitment = commitment
	s.commitment.ID.BatchID = s.batch.ID
	s.commitment.BodyHash = utils.NewRandomHash()

	s.batchNotFoundAPIErr = &APIError{
		Code:    30000,
		Message: "batch not found",
	}
}

func (s *GetBatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchTestSuite) TestGetBatchByHash() {
	s.addStateLeaf()
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*s.batch.Hash)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch.ID, result.ID)
	s.Equal(s.batch.Hash, result.Hash)
	s.Equal(s.batch.Type, result.Type)
	s.Equal(s.batch.TransactionHash, result.TransactionHash)
	s.Equal(*s.batch.FinalisationBlock-config.DefaultBlocksToFinalise, result.SubmissionBlock)
	s.Equal(s.batch.FinalisationBlock, result.FinalisationBlock)
	s.Equal(s.batch.SubmissionTime, result.SubmissionTime)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_GenesisBatch() {
	genesisBatch := models.Batch{
		ID:                models.MakeUint256(0),
		Type:              batchtype.Genesis,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
	}
	err := s.storage.AddBatch(&genesisBatch)
	s.NoError(err)

	result, err := s.api.GetBatchByHash(*genesisBatch.Hash)
	s.NoError(err)
	s.Equal(genesisBatch.Hash, result.Hash)
	s.Equal(genesisBatch.Type, result.Type)
	s.Equal(*genesisBatch.FinalisationBlock, result.SubmissionBlock)
	s.Len(result.Commitments, 0)
}

func (s *GetBatchTestSuite) TestGetBatchByHash_NonexistentBatch() {
	result, err := s.api.GetBatchByHash(utils.RandomHash())
	s.Equal(s.batchNotFoundAPIErr, err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) TestGetBatchByID() {
	s.addStateLeaf()
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	err = s.storage.AddTxCommitment(&s.commitment)
	s.NoError(err)

	result, err := s.api.GetBatchByID(s.batch.ID)
	s.NoError(err)
	s.NotNil(result)
	s.Len(result.Commitments, 1)
	s.Equal(s.batch.ID, result.ID)
	s.Equal(s.batch.Hash, result.Hash)
	s.Equal(s.batch.Type, result.Type)
	s.Equal(s.batch.TransactionHash, result.TransactionHash)
	s.Equal(*s.batch.FinalisationBlock-config.DefaultBlocksToFinalise, result.SubmissionBlock)
	s.Equal(s.batch.FinalisationBlock, result.FinalisationBlock)
	s.Equal(s.batch.SubmissionTime, result.SubmissionTime)
}

func (s *GetBatchTestSuite) TestGetBatchByID_GenesisBatch() {
	genesisBatch := models.Batch{
		ID:                models.MakeUint256(0),
		Type:              batchtype.Genesis,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		FinalisationBlock: ref.Uint32(10),
	}
	err := s.storage.AddBatch(&genesisBatch)
	s.NoError(err)

	result, err := s.api.GetBatchByID(genesisBatch.ID)
	s.NoError(err)
	s.Equal(genesisBatch.Hash, result.Hash)
	s.Equal(genesisBatch.Type, result.Type)
	s.Equal(*genesisBatch.FinalisationBlock, result.SubmissionBlock)
	s.Len(result.Commitments, 0)
}

func (s *GetBatchTestSuite) TestGetBatchByID_NonexistentBatch() {
	result, err := s.api.GetBatchByID(models.MakeUint256(0))
	s.Equal(s.batchNotFoundAPIErr, err)
	s.Nil(result)
}

func (s *GetBatchTestSuite) addStateLeaf() {
	_, err := s.storage.StateTree.Set(uint32(1), &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestGetBatchTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchTestSuite))
}
