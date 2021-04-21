package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NetworkInfoTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	db         *db.TestDB
	testClient *eth.TestClient
}

func (s *NetworkInfoTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *NetworkInfoTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)

	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	s.api = &API{nil, storage, s.testClient.Client}
	s.db = testDB
}

func (s *NetworkInfoTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *NetworkInfoTestSuite) TestGetNetworkInfo_NoBatches() {
	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal("", networkInfo.LatestBatch)
	s.Equal("", networkInfo.LatestFinalisedBatch)

}

func (s *NetworkInfoTestSuite) TestGetNetworkInfo_NoFinalisedBatches() {
	batches := []models.Batch{
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Transfer,
			ID:                models.MakeUint256(1234),
			FinalisationBlock: 1234,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			ID:                models.MakeUint256(2000),
			FinalisationBlock: 2000,
		},
	}
	err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal("2000", networkInfo.LatestBatch)
	s.Equal("", networkInfo.LatestFinalisedBatch)

}

func (s *NetworkInfoTestSuite) TestGetNetworkInfo() {
	batches := []models.Batch{
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Transfer,
			ID:                models.MakeUint256(1234),
			FinalisationBlock: 1,
		},
		{
			Hash:              utils.RandomHash(),
			Type:              txtype.Create2Transfer,
			ID:                models.MakeUint256(2000),
			FinalisationBlock: 2000,
		},
	}
	err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	s.api.storage.SetLatestBlockNumber(1)

	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal(uint32(1), networkInfo.BlockNumber)
	s.Equal("2000", networkInfo.LatestBatch)
	s.Equal("1234", networkInfo.LatestFinalisedBatch)
}

func TestNetworkInfoTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkInfoTestSuite))
}
