package mempool

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxPoolTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	txPool  *txPool
}

func (s *TxPoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxPoolTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.txPool, err = NewTxPool(s.storage.Storage)
	s.NoError(err)
}

func (s *TxPoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxPoolTestSuite) TestReadTxsAndUpdateMempool() {
	_, err := s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = s.txPool.ReadTxs(ctx)
		s.NoError(err)
	}()

	for i := 0; i < 5; i++ {
		s.txPool.Send(s.newTransfer(0, uint64(i)))
	}

	s.Eventually(func() bool {
		s.txPool.mutex.Lock()
		defer s.txPool.mutex.Unlock()
		return len(s.txPool.incomingTxs) == 5
	}, 1*time.Second, 10*time.Millisecond)

	err = s.txPool.UpdateMempool()
	s.NoError(err)

	receivedTxs := s.getAllTxs(0)
	s.Len(receivedTxs, 5)

	cancel()
	wg.Wait()
}

func (s *TxPoolTestSuite) newTransfer(from uint32, nonce uint64) *models.Transfer {
	return testutils.NewTransfer(from, 1, nonce, 100)
}

func (s *TxPoolTestSuite) getAllTxs(stateID uint32) []models.GenericTransaction {
	txs := s.txPool.Mempool().GetExecutableTxs(txtype.Transfer)

	_, txMempool := s.txPool.Mempool().BeginTransaction()
	for {
		tx := txMempool.GetNextExecutableTx(txtype.Transfer, stateID)
		if tx == nil {
			break
		}
		txs = append(txs, tx)
	}
	return txs
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(TxPoolTestSuite))
}
