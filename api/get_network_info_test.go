package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
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

	s.api = &API{storage: testStorage.InternalStorage, client: s.testClient.Client}
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
			ID:                models.MakeUint256(1234),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1234),
		},
		{
			ID:                models.MakeUint256(2000),
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2000),
		},
	}
	err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	networkInfo, err := s.api.GetNetworkInfo()
	s.NoError(err)
	s.NotNil(networkInfo)
	s.Equal("2000", networkInfo.LatestBatch.String())
	s.Nil(networkInfo.LatestFinalisedBatch)
}

func (s *NetworkInfoTestSuite) TestGetNetworkInfo() {
	batches := []models.Batch{
		{
			ID:                models.MakeUint256(1234),
			Type:              txtype.Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(1),
		},
		{
			ID:                models.MakeUint256(2000),
			Type:              txtype.Create2Transfer,
			TransactionHash:   utils.RandomHash(),
			Hash:              utils.NewRandomHash(),
			FinalisationBlock: ref.Uint32(2000),
		},
	}
	err := s.api.storage.AddBatch(&batches[0])
	s.NoError(err)
	err = s.api.storage.AddBatch(&batches[1])
	s.NoError(err)

	commitmentInBatch := commitment
	commitmentInBatch.IncludedInBatch = &batches[0].ID
	commitmentID, err := s.api.storage.AddCommitment(&commitmentInBatch)
	s.NoError(err)
	_, err = s.api.storage.AddTransfer(&models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:                 common.Hash{1, 2, 3},
			TxType:               txtype.Transfer,
			FromStateID:          0,
			IncludedInCommitment: commitmentID,
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
	s.Equal("2000", networkInfo.LatestBatch.String())
	s.Equal("1234", networkInfo.LatestFinalisedBatch.String())
	s.Equal(expectedDomain, networkInfo.SignatureDomain.Bytes())
}

func TestNetworkInfoTestSuite(t *testing.T) {
	suite.Run(t, new(NetworkInfoTestSuite))
}
