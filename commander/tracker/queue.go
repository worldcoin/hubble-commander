package tracker

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
)

type txsQueue struct {
	txs   []*types.Transaction
	mutex sync.RWMutex
}

func newTxsQueue() *txsQueue {
	return &txsQueue{
		txs: make([]*types.Transaction, 0),
	}
}

func (q *txsQueue) Add(tx *types.Transaction) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.txs = append(q.txs, tx)
}

func (q *txsQueue) First() *types.Transaction {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	if len(q.txs) == 0 {
		return nil
	}
	return q.txs[0]
}

func (q *txsQueue) RemoveFirst() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.txs = q.txs[1:]
}

func (q *txsQueue) IsEmpty() bool {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.txs) == 0
}
