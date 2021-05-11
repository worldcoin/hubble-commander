package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
	batch   models.Batch
}

func (s *GetBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	ethClient, err := eth.NewTestClient()
	s.NoError(err)
	s.api = &API{nil, s.storage.Storage, ethClient.Client}

	s.batch = models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
		FinalisationBlock: 42000,
	}
}

func (s *GetBatchesTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	err := s.storage.AddBatch(&s.batch)
	s.NoError(err)

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 1)
	s.Equal(s.batch, result[0].Batch)
	s.Equal(s.batch.FinalisationBlock-rollup.DefaultBlocksToFinalise, result[0].SubmissionBlock)
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
