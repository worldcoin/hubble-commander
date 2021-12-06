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

type GetBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	storage    *st.TestStorage
	testClient *eth.TestClient
	batches    []models.Batch
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

	s.batches = []models.Batch{
		{
			ID:                models.MakeUint256(1),
			Type:              batchtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(42000),
			SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
		},
		{
			ID:                models.MakeUint256(2),
			Type:              batchtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(43000),
			SubmissionTime:    models.NewTimestamp(time.Unix(150, 0).UTC()),
		},
		{
			ID:                models.MakeUint256(3),
			Type:              batchtype.MassMigration,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(44000),
			SubmissionTime:    models.NewTimestamp(time.Unix(160, 0).UTC()),
		},
	}
}

func (s *GetBatchesTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetBatchesTestSuite) TestGetBatches() {
	s.addBatches()

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(3))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 3)

	for i := range s.batches {
		s.Equal(s.batches[i].ID, result[i].ID)
		s.Equal(s.batches[i].Hash, result[i].Hash)
		s.Equal(s.batches[i].Type, result[i].Type)
		s.Equal(s.batches[i].TransactionHash, result[i].TransactionHash)
		s.Equal(*s.batches[i].FinalisationBlock-config.DefaultBlocksToFinalise, result[i].SubmissionBlock)
		s.Equal(s.batches[i].FinalisationBlock, result[i].FinalisationBlock)
		s.Equal(s.batches[i].SubmissionTime, result[i].SubmissionTime)
	}
}

func (s *GetBatchesTestSuite) TestGetBatches_PendingBatch() {
	pendingBatch := s.batches[0]
	pendingBatch.Hash = nil
	pendingBatch.FinalisationBlock = nil
	err := s.storage.AddBatch(&pendingBatch)
	s.NoError(err)

	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 0)
}

func (s *GetBatchesTestSuite) TestGetBatchesByHash_NoBatches() {
	result, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1))
	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 0)
}

func (s *GetBatchesTestSuite) addBatches() {
	for i := range s.batches {
		err := s.storage.AddBatch(&s.batches[i])
		s.NoError(err)
	}
}

func TestGetBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetBatchesTestSuite))
}
