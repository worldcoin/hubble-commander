package admin

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

const authKeyValue = "secret key"

type GetPendingBatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
	client  *eth.TestClient
	batches []models.Batch
}

func (s *GetPendingBatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetPendingBatchesTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.api = &API{
		cfg:     &config.APIConfig{AuthenticationKey: ref.String(authKeyValue)},
		storage: s.storage.Storage,
		client:  s.client.Client,
	}

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
			ID:              models.MakeUint256(2),
			Type:            batchtype.Create2Transfer,
			TransactionHash: utils.RandomHash(),
		},
		{
			ID:              models.MakeUint256(3),
			Type:            batchtype.MassMigration,
			TransactionHash: utils.RandomHash(),
		},
	}
}

func (s *GetPendingBatchesTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches() {
	s.addBatches()

	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)

	expected := s.batches[1:]
	for i := range batches {
		s.Equal(expected[i].ID, batches[i].ID)
		s.Equal(expected[i].Type, batches[i].Type)
		s.Equal(expected[i].TransactionHash, batches[i].TransactionHash)
	}
}

func (s *GetPendingBatchesTestSuite) TestGetPendingBatches_NoBatches() {
	batches, err := s.api.GetPendingBatches(contextWithAuthKey(authKeyValue))
	s.NoError(err)
	s.Len(batches, 0)
}

func (s *GetPendingBatchesTestSuite) addBatches() {
	for i := range s.batches {
		err := s.storage.AddBatch(&s.batches[i])
		s.NoError(err)
	}
}

func TestGetPendingBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(GetPendingBatchesTestSuite))
}
