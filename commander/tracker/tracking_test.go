package tracker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxsTrackingTestSuite struct {
	*require.Assertions
	suite.Suite
	client            *eth.TestClient
	txsChannels       *eth.TxsTrackingChannels
	wg                sync.WaitGroup
	cancelTxsTracking context.CancelFunc
	txsQueue          *txsQueue
}

func (s *TxsTrackingTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxsTrackingTestSuite) SetupTest() {
	s.wg = sync.WaitGroup{}
	s.txsChannels = &eth.TxsTrackingChannels{
		SkipSendingRequestsThroughChannel: true,
		SentTxs:                           make(chan *types.Transaction, 8),
	}

	var err error
	s.client, err = eth.NewConfiguredTestClient(
		&rollup.DeploymentConfig{},
		&eth.TestClientConfig{
			TxsChannels: s.txsChannels,
		},
	)
	s.NoError(err)
	s.txsQueue = newTxsQueue()
	s.startTxsTracking()
}

func (s *TxsTrackingTestSuite) startTxsTracking() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelTxsTracking = cancel

	s.wg.Add(1)
	go func() {
		err := startTrackingSentTxs(ctx, s.client.Client, s.txsChannels.SentTxs, s.txsQueue)
		s.NoError(err)
		s.wg.Done()
	}()
}

func (s *TxsTrackingTestSuite) TearDownTest() {
	s.cancelTxsTracking()
	s.wg.Wait()
	s.client.Close()
}

func (s *TxsTrackingTestSuite) TestStartTrackingSentTxs_ChannelBufferOverflowWithoutBlocking() {
	txs := make([]*types.Transaction, 20)
	commitments := getCommitments(batchtype.Transfer)

	for i := 0; i < len(txs); i++ {
		tx, err := s.client.SubmitTransfersBatch(models.NewUint256(uint64(i+1)), commitments)
		s.NoError(err)
		txs[i] = tx
	}

	for _, tx := range txs {
		s.txsChannels.SentTxs <- tx
	}

	s.Eventually(func() bool {
		return s.txsQueue.IsEmpty()
	}, time.Second, time.Millisecond*300)
}

func TestTxsTrackingTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackingTestSuite))
}
