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
	"github.com/Worldcoin/hubble-commander/utils/ref"
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

	setUserStates(s.Assertions, s.storage.StateTree, map[uint32]uint64{
		0: 5,
		1: 0,
	})
}

func (s *TxPoolTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxPoolTestSuite) TestNewTxPool_HandlesValidReplacements() {
	txs := []models.GenericTransaction{
		s.newTransferAt(1, 0, 5, 100), // 0 - ok
		s.newC2TAt(2, 0, 4, 100),      // 1 - nonce too low

		s.newTransferAt(3, 1, 0, 100), // 2
		s.newC2TAt(4, 1, 0, 200),      // 3 - valid replacement of 2
		s.newTransferAt(5, 1, 0, 150), // 4 - invalid replacement of 3
	}
	// no need to shuffle txs as they are retrieved from DB sorted by tx hashes which are random
	err := s.storage.BatchAddTransaction(models.MakeGenericArray(txs...))
	s.NoError(err)

	s.txPool, err = NewTxPool(s.storage.Storage)
	s.NoError(err)

	mempoolTxs := s.getAllMempoolTxs([]uint32{0, 1})
	s.Len(mempoolTxs, 2)
	s.Contains(mempoolTxs, txs[0])
	s.Contains(mempoolTxs, txs[3])

	tx1, err := s.storage.GetCreate2Transfer(txs[1].GetBase().Hash)
	s.NoError(err)
	s.Equal(ErrTxNonceTooLow.Error(), *tx1.ErrorMessage)

	_, err = s.storage.GetTransfer(txs[2].GetBase().Hash)
	s.True(st.IsNotFoundError(err))

	tx4, err := s.storage.GetTransfer(txs[4].GetBase().Hash)
	s.NoError(err)
	s.Equal(ErrTxReplacementFailed.Error(), *tx4.ErrorMessage)
}

func (s *TxPoolTestSuite) TestReadTxsAndUpdateMempool() {
	stopReadingTxs := s.startReadingTxs()
	defer stopReadingTxs()

	for i := 5; i < 10; i++ {
		s.txPool.Send(s.newTransfer(uint64(i), 10))
	}

	s.waitForTxsToBeRead(5)

	err := s.txPool.UpdateMempool()
	s.NoError(err)

	receivedTxs := s.getAllTransfers()
	s.Len(receivedTxs, 5)
}

func (s *TxPoolTestSuite) TestUpdateMempool_MarksInvalidReplacementTxAsFailed() {
	tx := s.newTransfer(5, 10)
	replacementTx := s.newTransfer(5, 5)
	err := s.storage.AddTransaction(replacementTx)
	s.NoError(err)

	stopReadingTxs := s.startReadingTxs()
	defer stopReadingTxs()

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

	mempoolTxs := s.getAllTransfers()
	s.Len(mempoolTxs, 1)
	s.Equal(tx, mempoolTxs[0])
}

func (s *TxPoolTestSuite) TestUpdateMempool_ReplacesPendingTx() {
	previousTx := s.newTransfer(5, 5)
	newTx := s.newTransfer(5, 10)

	err := s.storage.AddTransaction(previousTx)
	s.NoError(err)

	stopReadingTxs := s.startReadingTxs()
	defer stopReadingTxs()

	s.txPool.Send(previousTx)
	s.txPool.Send(newTx)

	s.waitForTxsToBeRead(2)

	err = s.txPool.UpdateMempool()
	s.NoError(err)

	_, err = s.storage.GetTransfer(previousTx.Hash)
	s.True(st.IsNotFoundError(err))

	mempoolTxs := s.getAllTransfers()
	s.Len(mempoolTxs, 1)
	s.Equal(newTx, mempoolTxs[0])
}

func (s *TxPoolTestSuite) TestUpdateMempool_RemovesPendingTxsWithTooLowNonces() {
	invalidTxs := []models.GenericTransaction{
		s.newTransfer(0, 10),
		s.newTransfer(1, 10),
	}
	validTx := s.newTransfer(5, 10)
	txs := []models.GenericTransaction{
		invalidTxs[0],
		validTx,
		invalidTxs[1],
	}

	stopReadingTxs := s.startReadingTxs()
	defer stopReadingTxs()

	for _, tx := range txs {
		err := s.storage.AddTransaction(tx)
		s.NoError(err)
		s.txPool.Send(tx)
	}

	s.waitForTxsToBeRead(3)

	err := s.txPool.UpdateMempool()
	s.NoError(err)

	mempoolTxs := s.getAllTransfers()
	s.Len(mempoolTxs, 1)
	s.Equal(validTx, mempoolTxs[0])

	failedTxs, err := s.storage.GetAllFailedTransactions()
	s.NoError(err)
	s.Len(failedTxs, 2)
	for _, badTx := range invalidTxs {
		badTx.GetBase().ErrorMessage = ref.String(ErrTxNonceTooLow.Error())
		s.Contains(failedTxs, badTx)
	}
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

func (s *TxPoolTestSuite) startReadingTxs() (stopReadingTxs func()) {
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.txPool.ReadTxs(ctx)
		s.NoError(err)
	}()
	return func() {
		cancel()
		wg.Wait()
	}
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

func (s *TxPoolTestSuite) newTransferAt(timestamp int64, from uint32, nonce, fee uint64) *models.Transfer {
	tx := testutils.NewTransfer(from, 1, nonce, 100)
	tx.GetBase().Fee = models.MakeUint256(fee)
	tx.GetBase().ReceiveTime = models.NewTimestamp(time.Unix(timestamp, 0).UTC())
	return tx
}

func (s *TxPoolTestSuite) newC2TAt(timestamp int64, from uint32, nonce, fee uint64) *models.Create2Transfer {
	tx := testutils.NewCreate2Transfer(from, ref.Uint32(1), nonce, 100, nil)
	tx.GetBase().Fee = models.MakeUint256(fee)
	tx.GetBase().ReceiveTime = models.NewTimestamp(time.Unix(timestamp, 0).UTC())
	return tx
}

func (s *TxPoolTestSuite) getAllTransfers() []models.GenericTransaction {
	return s.getAllUsersTxs(txtype.Transfer, []uint32{0})
}

func (s *TxPoolTestSuite) getAllMempoolTxs(stateIDs []uint32) []models.GenericTransaction {
	txs := make([]models.GenericTransaction, 0)
	for txType := range txtype.TransactionTypes {
		userTxs := s.getAllUsersTxs(txType, stateIDs)
		txs = append(txs, userTxs...)
	}
	return txs
}

func (s *TxPoolTestSuite) getAllUsersTxs(txType txtype.TransactionType, stateIDs []uint32) []models.GenericTransaction {
	txs := s.txPool.Mempool().GetExecutableTxs(txType)

	_, txMempool := s.txPool.Mempool().BeginTransaction()
	for _, stateID := range stateIDs {
		for {
			tx, err := txMempool.GetNextExecutableTx(txType, stateID)
			s.NoError(err)
			if tx == nil {
				break
			}
			txs = append(txs, tx)
		}
	}
	return txs
}

func TestPoolTestSuite(t *testing.T) {
	suite.Run(t, new(TxPoolTestSuite))
}
