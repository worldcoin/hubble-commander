package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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
	}
	s.tokenID = models.NewUint256(0) // First registered tokenID
}

func (s *DepositsTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestSyncQueuedDeposits() {
	s.registerToken()
	s.approveTokens()

	deposit := s.queueDeposit()

	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	queuedDeposits, err := s.cmd.syncQueuedDeposits(0, *latestBlockNumber)
	s.NoError(err)

	s.Len(queuedDeposits, 1)
	s.Contains(queuedDeposits, *deposit)
}

func (s *DepositsTestSuite) registerToken() {
	token := models.RegisteredToken{
		Contract: s.testClient.ExampleTokenAddress,
	}
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)
	RegisterSingleToken(s.Assertions, s.testClient, &token, latestBlockNumber)
}

func (s *DepositsTestSuite) approveTokens() {
	token, err := erc20.NewERC20(s.testClient.ExampleTokenAddress, s.testClient.GetBackend())
	s.NoError(err)

	tx, err := token.Approve(s.testClient.GetAccount(), s.testClient.ChainState.DepositManager, utils.ParseEther("10"))
	s.NoError(err)

	_, err = deployer.WaitToBeMined(s.testClient.GetBackend(), tx)
	s.NoError(err)
}

func (s *DepositsTestSuite) queueDeposit() *models.Deposit {
	deposits, unsubscribe, err := s.testClient.WatchQueuedDeposits(&bind.WatchOpts{Start: nil})
	s.NoError(err)
	defer unsubscribe()

	toPubKeyID := models.NewUint256(1)
	l1Amount := models.NewUint256FromBig(*utils.ParseEther("10"))
	depositID, l2Amount, err := s.testClient.QueueDeposit(toPubKeyID, l1Amount, s.tokenID, deposits)
	s.NoError(err)

	return &models.Deposit{
		ID:                   *depositID,
		ToPubKeyID:           uint32(toPubKeyID.Uint64()),
		TokenID:              *s.tokenID,
		L2Amount:             *l2Amount,
		IncludedInCommitment: nil,
	}
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
