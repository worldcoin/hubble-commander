package tracker

import (
	"sync"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsTrackerTestSuite struct {
	*require.Assertions
	suite.Suite
	TestSuiteWithTxsSending
	testClient *eth.TestClient
}

func (s *TxsTrackerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxsTrackerTestSuite) SetupTest() {
	var err error
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	s.StartTxsSending(s.testClient.TxsChannels.Requests)
}

func (s *TxsTrackerTestSuite) TearDownTest() {
	s.StopTxsSending()
	s.testClient.Close()
}

func (s *TxsTrackerTestSuite) TestTxsTracker_SendTransactionsAtTheSameTime() {
	start := make(chan struct{})
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)
	resultTxs := make([]*types.Transaction, 2)
	batchID := models.NewUint256(1)

	go func() {
		var err error
		<-start
		resultTxs[0], err = s.testClient.WithdrawStake(batchID)
		s.NoError(err)
		waitGroup.Done()
	}()

	go func() {
		var err error
		commitments := getCommitments(batchtype.Transfer)
		<-start
		resultTxs[1], err = s.testClient.SubmitTransfersBatch(batchID, commitments)
		s.NoError(err)
		waitGroup.Done()
	}()

	close(start)
	waitGroup.Wait()

	s.NotEqual(resultTxs[0].Nonce(), resultTxs[1].Nonce())
}

func getCommitments(batchType batchtype.BatchType) []models.CommitmentWithTxs {
	commitment := models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          batchType,
				PostStateRoot: utils.RandomHash(),
			},
			FeeReceiver:       uint32(1234),
			CombinedSignature: models.MakeRandomSignature(),
		},
		Transactions: utils.RandomBytes(12),
	}
	return []models.CommitmentWithTxs{&commitment}
}

func TestTxsTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackerTestSuite))
}
