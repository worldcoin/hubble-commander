package tracker

import (
	"sync"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
)

type Tracker struct {
	txs   []*types.Transaction
	mutex sync.RWMutex

	client       *eth.Client
	txsChan      <-chan *types.Transaction
	requestsChan <-chan *eth.TxSendingRequest
}

func NewTracker(client *eth.Client, txsChan <-chan *types.Transaction, requestsChan <-chan *eth.TxSendingRequest) *Tracker {
	return &Tracker{
		txs:          make([]*types.Transaction, 0),
		client:       client,
		txsChan:      txsChan,
		requestsChan: requestsChan,
	}
}

func (t *Tracker) addTx(tx *types.Transaction) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.txs = append(t.txs, tx)
}

func (t *Tracker) firstTx() *types.Transaction {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	if len(t.txs) == 0 {
		return nil
	}
	return t.txs[0]
}

func (t *Tracker) removeFirstTx() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.txs = t.txs[1:]
}

func (t *Tracker) isEmptyTxsQueue() bool {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return len(t.txs) == 0
}
