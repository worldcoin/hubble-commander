package tracker

import (
	"context"
	"sync"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsSendingTestSuite struct {
	*require.Assertions
	suite.Suite
	client           *eth.TestClient
	txsChannels      *eth.TxsTrackingChannels
	wg               sync.WaitGroup
	cancelTxsSending context.CancelFunc
}

func (s *TxsSendingTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxsSendingTestSuite) SetupTest() {
	s.wg = sync.WaitGroup{}
	s.txsChannels = &eth.TxsTrackingChannels{
		Requests: make(chan *eth.TxSendingRequest, 8),
		SentTxs:  make(chan *types.Transaction, 8),
	}

	var err error
	s.client, err = eth.NewConfiguredTestClient(
		&rollup.DeploymentConfig{},
		&eth.TestClientConfig{
			TxsChannels: s.txsChannels,
		},
	)
	s.NoError(err)
	s.startTxsSending()
}

func (s *TxsSendingTestSuite) startTxsSending() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelTxsSending = cancel

	s.wg.Add(1)
	go func() {
		err := SendRequestedTxs(ctx, s.txsChannels.Requests)
		s.NoError(err)
		s.wg.Done()
	}()
}

func (s *TxsSendingTestSuite) TearDownTest() {
	s.cancelTxsSending()
	s.wg.Wait()
	s.client.Close()
}

func (s *TxsSendingTestSuite) TestSendRequestedTxs_SetsConsecutiveNoncesForTxsSentInSameTime() {
	start := make(chan struct{})
	waitGroup := sync.WaitGroup{}
	resultTxs := make([]*types.Transaction, 2)
	batchID := models.NewUint256(1)

	waitGroup.Add(1)
	go func() {
		var err error
		<-start
		resultTxs[0], err = s.client.WithdrawStake(batchID)
		s.NoError(err)
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		var err error
		commitments := getCommitments(batchtype.Transfer)
		<-start
		resultTxs[1], err = s.client.SubmitTransfersBatch(batchID, commitments)
		s.NoError(err)
		waitGroup.Done()
	}()

	close(start)
	waitGroup.Wait()

	s.NotEqual(resultTxs[0].Nonce(), resultTxs[1].Nonce())
}

func getCommitments(batchType batchtype.BatchType) []models.CommitmentWithTxs {
	return []models.CommitmentWithTxs{
		&models.TxCommitmentWithTxs{
			TxCommitment: models.TxCommitment{
				CommitmentBase: models.CommitmentBase{
					Type:          batchType,
					PostStateRoot: utils.RandomHash(),
				},
				FeeReceiver:       uint32(1234),
				CombinedSignature: models.MakeRandomSignature(),
			},
			Transactions: utils.RandomBytes(12),
		},
	}
}

func TestTxsSendingTestSuite(t *testing.T) {
	suite.Run(t, new(TxsSendingTestSuite))
}
