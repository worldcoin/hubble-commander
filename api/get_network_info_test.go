package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type NetworkInfoTestSuite struct {
	*require.Assertions
	suite.Suite
	api        *API
	teardown   func() error
	testClient *eth.TestClient
}

func (s *NetworkInfoTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *NetworkInfoTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	err = testStorage.SetChainState(&chainState)
	s.NoError(err)

	s.api = &API{storage: testStorage.Storage, client: s.testClient.Client}
}

func (s *NetworkInfoTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
	s.testClient.Close()
}

func (s *NetworkInfoTestSuite) TestGetNetworkInfo_NoBatches() {
	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Nil(networkInfo.LatestBatch)
	s.Nil(networkInfo.LatestFinalisedBatch)
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
	_, err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	_, err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal("2000", *networkInfo.LatestBatch)
	s.Nil(networkInfo.LatestFinalisedBatch)
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
	_, err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	_, err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	err = s.api.storage.AddTransfer(&models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.Hash{1, 2, 3},
			TxType:      txtype.Transfer,
			FromStateID: 0,
		},
		ToStateID: 1,
	})
	s.NoError(err)

	s.api.storage.SetLatestBlockNumber(1)
	expectedDomain := crypto.Keccak256(chainState.Rollup.Bytes())

	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal(uint32(1), networkInfo.BlockNumber)
	s.Equal(uint32(0), networkInfo.AccountCount)
	s.Equal(1, networkInfo.TransactionCount)
	s.Equal("2000", *networkInfo.LatestBatch)
	s.Equal("1234", *networkInfo.LatestFinalisedBatch)
	s.Equal(expectedDomain, networkInfo.SignatureDomain.Bytes())
}

func TestNetworkInfoTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkInfoTestSuite))
}
