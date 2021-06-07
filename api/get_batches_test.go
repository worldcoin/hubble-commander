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

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.TestStorage
	testClient *eth.TestClient
	batch      models.Batch
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage, client: s.testClient.Client}

	s.batch = models.Batch{
		ID:                1,
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              utils.NewRandomHash(),
		Number:            models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(42000),
	}
}

func (s *GetBatchesTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	_, err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 1)
	s.Equal(s.batch.Number, result[0].ID)
	s.Equal(s.batch.Hash, result[0].Hash)
	s.Equal(s.batch.Type, result[0].Type)
	s.Equal(s.batch.TransactionHash, result[0].TransactionHash)
	s.Equal(*s.batch.FinalisationBlock-rollup.DefaultBlocksToFinalise, result[0].SubmissionBlock)
	s.Equal(s.batch.FinalisationBlock, result[0].FinalisationBlock)
}

func (s *GetBatchesTestSuite) TestGetBatchesByHash_NoBatches() {
	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 0)
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
