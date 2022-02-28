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

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(100),
		Nonce:    models.MakeUint256(5),
	})
	s.NoError(err)
}

func (s *TxPoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxPoolTestSuite) TestReadTxsAndUpdateMempool() {
	wg, cancel := s.startReadingTxs()

	for i := 5; i < 10; i++ {
		s.txPool.Send(s.newTransfer(uint64(i), 10))
	}

	s.waitForTxsToBeRead(5)

	err := s.txPool.UpdateMempool()
	s.NoError(err)

	receivedTxs := s.getAllTxs(0)
	s.Len(receivedTxs, 5)

	cancel()
	wg.Wait()
}

func (s *TxPoolTestSuite) TestUpdateMempool_MarksInvalidReplacementTxAsFailed() {
	tx := s.newTransfer(5, 10)
	replacementTx := s.newTransfer(5, 5)
	err := s.storage.AddTransaction(replacementTx)
	s.NoError(err)

	wg, cancel := s.startReadingTxs()

	s.txPool.Send(tx)
	s.txPool.Send(replacementTx)

	s.waitForTxsToBeRead(2)

	err = s.txPool.UpdateMempool()
	s.NoError(err)

	txs, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(txs, 1)
	s.Equal(replacementTx.Hash, txs.At(0).GetBase().Hash)
	s.Equal(ErrTxReplacementFailed.Error(), *txs.At(0).GetBase().ErrorMessage)

	mempoolTxs := s.getAllTxs(0)
	s.Len(mempoolTxs, 1)
	s.Equal(tx, mempoolTxs[0])

	cancel()
	wg.Wait()
}

func (s *TxPoolTestSuite) TestUpdateMempool_ReplacesPendingTx() {
	previousTx := s.newTransfer(5, 5)
	newTx := s.newTransfer(5, 10)

	err := s.storage.AddTransaction(previousTx)
	s.NoError(err)

	wg, cancel := s.startReadingTxs()

	s.txPool.Send(previousTx)
	s.txPool.Send(newTx)

	s.waitForTxsToBeRead(2)

	err = s.txPool.UpdateMempool()
	s.NoError(err)

	_, err = s.storage.GetTransfer(previousTx.Hash)
	s.True(st.IsNotFoundError(err))

	mempoolTxs := s.getAllTxs(0)
	s.Len(mempoolTxs, 1)
	s.Equal(newTx, mempoolTxs[0])

	cancel()
	wg.Wait()
}

func (s *TxPoolTestSuite) TestUpdateMempool_RemovesPendingTransactionsWithTooLowNonces() {

}

func (s *TxPoolTestSuite) TestRemoveFailedTxs_RemovesTxsFromMempoolAndMarksTxsAsFailed() {
	txs := []models.GenericTransaction{
		s.newTransfer(5, 10),
		s.newTransfer(6, 10),
		s.newTransfer(7, 10),
	}
	for _, tx := range txs {
		err := s.storage.AddTransaction(tx)
		s.NoError(err)
		s.txPool.addIncomingTx(tx)
	}

	err := s.txPool.UpdateMempool()
	s.NoError(err)

	err = s.txPool.RemoveFailedTxs(txsToTxErrors(txs...))
	s.NoError(err)

	failedTxs, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(failedTxs, 3)

	s.NotContains(s.txPool.mempool.buckets, 0)
}

func (s *TxPoolTestSuite) startReadingTxs() (*sync.WaitGroup, context.CancelFunc) {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.txPool.ReadTxs(ctx)
		s.NoError(err)
	}()
	return wg, cancel
}

func (s *TxPoolTestSuite) waitForTxsToBeRead(expectedTxsLength int) {
	s.Eventually(func() bool {
		s.txPool.mutex.Lock()
		defer s.txPool.mutex.Unlock()
		return len(s.txPool.incomingTxs) == expectedTxsLength
	}, 1*time.Second, 10*time.Millisecond)
}

func (s *TxPoolTestSuite) newTransfer(nonce, fee uint64) *models.Transfer {
	tx := testutils.NewTransfer(0, 1, nonce, 100)
	tx.Fee = models.MakeUint256(fee)
	return tx
}

func (s *TxPoolTestSuite) getAllTxs(stateID uint32) []models.GenericTransaction {
	txs := s.txPool.Mempool().GetExecutableTxs(txtype.Transfer)

	_, txMempool := s.txPool.Mempool().BeginTransaction()
	for {
		tx, err := txMempool.GetNextExecutableTx(txtype.Transfer, stateID)
		s.NoError(err)
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
