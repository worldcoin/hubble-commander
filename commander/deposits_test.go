package commander

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown   func() error
	testClient *eth.TestClient
	cmd        *Commander
	tokenID    *models.Uint256
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositsTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	s.cmd = &Commander{
		storage: testStorage.Storage,
		client:  s.testClient.Client,
		cfg:     &config.Config{},
		metrics: metrics.NewCommanderMetrics(),
	}
	s.tokenID = models.NewUint256(0) // First registered tokenID
}

func (s *DepositsTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestSyncDeposits() {
	s.approveTokens()

	// Smart contract needs 4 deposits to create a subtree (depth specified in rollup contract deployment)
	deposits := s.queueFourDeposits()

	s.queueDeposit()
	s.queueDeposit()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.syncDeposits(context.Background(), 0, *latestBlockNumber)
	s.NoError(err)

	subtree, err := s.cmd.storage.GetPendingDepositSubtree(models.MakeUint256(1))
	s.NoError(err)

	_, err = s.cmd.storage.GetPendingDepositSubtree(models.MakeUint256(2))
	s.True(st.IsNotFoundError(err))

	s.Equal(deposits, subtree.Deposits)

	_, err = s.cmd.storage.GetFirstPendingDeposits(4)
	s.ErrorIs(err, st.ErrRanOutOfPendingDeposits)
}

func (s *DepositsTestSuite) TestSyncDeposits_TwoSubtrees() {
	s.approveTokens()

	firstSubtreeDeposits := s.queueFourDeposits()
	secondSubtreeDeposits := s.queueFourDeposits()
	s.queueDeposit()
	s.queueDeposit()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.syncDeposits(context.Background(), 0, *latestBlockNumber)
	s.NoError(err)

	firstSubtree, err := s.cmd.storage.GetPendingDepositSubtree(models.MakeUint256(1))
	s.NoError(err)

	secondSubtree, err := s.cmd.storage.GetPendingDepositSubtree(models.MakeUint256(2))
	s.NoError(err)

	_, err = s.cmd.storage.GetPendingDepositSubtree(models.MakeUint256(3))
	s.True(st.IsNotFoundError(err))

	s.Equal(firstSubtreeDeposits, firstSubtree.Deposits)
	s.Equal(secondSubtreeDeposits, secondSubtree.Deposits)

	_, err = s.cmd.storage.GetFirstPendingDeposits(4)
	s.ErrorIs(err, st.ErrRanOutOfPendingDeposits)
}

func (s *DepositsTestSuite) TestSyncQueuedDeposits() {
	s.approveTokens()

	deposit := s.queueDeposit()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.syncQueuedDeposits(0, *latestBlockNumber)
	s.NoError(err)

	syncedDeposits, err := s.cmd.storage.GetFirstPendingDeposits(1)
	s.NoError(err)
	s.Equal(*deposit, syncedDeposits[0])
}

func (s *DepositsTestSuite) TestFetchDepositSubtrees() {
	s.approveTokens()

	// Smart contract needs 4 deposits to create a subtree (depth specified in cfg.Rollup.MaxDepositSubtreeDepth)
	s.queueFourDeposits()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	depositSubtrees, err := s.cmd.fetchDepositSubtrees(0, *latestBlockNumber)
	s.NoError(err)

	s.Len(depositSubtrees, 1)
	s.Equal(depositSubtrees[0].ID, models.MakeUint256(1))
	s.NotEqual(depositSubtrees[0].Root, common.Hash{})
	s.Nil(depositSubtrees[0].Deposits)
}

func (s *DepositsTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.testClient.ExampleTokenAddress, s.testClient.GetBackend())
	s.NoError(err)

	tx, err := token.Approve(s.testClient.GetAccount(), s.testClient.ChainState.DepositManager, utils.ParseEther("100"))
	s.NoError(err)

	_, err = s.testClient.WaitToBeMined(tx)
	s.NoError(err)
}

func (s *DepositsTestSuite) queueDeposit() *models.PendingDeposit {
	toPubKeyID := models.NewUint256(1)
	l1Amount := models.NewUint256FromBig(*utils.ParseEther("10"))
	depositID, l2Amount, err := s.testClient.QueueDepositAndWait(toPubKeyID, l1Amount, s.tokenID)
	s.NoError(err)

	return &models.PendingDeposit{
		ID:         *depositID,
		ToPubKeyID: uint32(toPubKeyID.Uint64()),
		TokenID:    *s.tokenID,
		L2Amount:   *l2Amount,
	}
}

func (s *DepositsTestSuite) queueFourDeposits() []models.PendingDeposit {
	return []models.PendingDeposit{
		*s.queueDeposit(),
		*s.queueDeposit(),
		*s.queueDeposit(),
		*s.queueDeposit(),
	}
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
