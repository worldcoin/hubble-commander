package mempool

import (
	"context"
	"sort"
	"sync"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type TxPool interface {
	Send(tx models.GenericTransaction)
	ReadTxs(ctx context.Context) error
	UpdateMempool() error
	RemoveFailedTxs(txErrors []models.TxError) error
	Mempool() *Mempool
}

type txPool struct {
	storage         *st.Storage
	mempool         *Mempool
	incomingTxs     []models.GenericTransaction
	incomingTxsChan chan models.GenericTransaction

	mutex sync.Mutex
}

func NewTxPoolWithGauge(storage *st.Storage, sizeGauge prometheus.Gauge) (*txPool, error) {
	txPool, err := NewTxPool(storage)
	if err != nil {
		return txPool, err
	}

	if sizeGauge != nil {
		sizeGauge.Set(0)
	}

	txPool.mempool.sizeGauge = sizeGauge
	return txPool, nil
}

func NewTxPool(storage *st.Storage) (*txPool, error) {
	txs, err := storage.GetAllPendingTransactions()
	if err != nil {
		return nil, err
	}

	sort.Slice(txs, func(i, j int) bool {
		return earlierTimestamp(txs[i].GetBase().ReceiveTime, txs[j].GetBase().ReceiveTime)
	})

	pool := txPool{
		storage:         storage,
		mempool:         NewMempool(),
		incomingTxs:     txs.ToSlice(),
		incomingTxsChan: make(chan models.GenericTransaction, 1024),
	}

	err = pool.UpdateMempool()
	if err != nil {
		return nil, err
	}

	return &pool, nil
}
func (p *txPool) ReadTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case tx := <-p.incomingTxsChan:
			p.addIncomingTx(tx)
		}
	}
}

func (p *txPool) addIncomingTx(tx models.GenericTransaction) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.incomingTxs = append(p.incomingTxs, tx)
}

func (p *txPool) UpdateMempool() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.incomingTxs) == 0 {
		return nil
	}

	for _, tx := range p.incomingTxs {
		err := p.addOrReplaceTx(tx)
		if err != nil {
			return err
		}
	}

	p.incomingTxs = make([]models.GenericTransaction, 0)
	return nil
}

func (p *txPool) addOrReplaceTx(tx models.GenericTransaction) error {
	txHash := &tx.GetBase().Hash
	prevTxHash, err := p.mempool.AddOrReplace(p.storage, tx)
	if errors.Is(err, ErrTxNonceTooLow) {
		log.WithField("txHash", *txHash).Errorf("Mempool: %s failed: %s", tx.Type().String(), err)
		return p.storage.SetTransactionError(getNonceError(txHash))
	}
	if errors.Is(err, ErrTxReplacementFailed) {
		log.WithField("txHash", *txHash).Debug("Mempool: transaction replacement failed")
		return p.storage.SetTransactionError(getReplacementError(txHash))
	}
	if err != nil {
		return err
	}

	if prevTxHash != nil {
		log.WithFields(log.Fields{
			"previousTxHash": *prevTxHash,
			"newTxHash":      *txHash,
		}).Debug("Mempool: replaced transaction")
		err = p.storage.RemovePendingTransactions(*prevTxHash)
		if st.IsNotFoundError(err) {
			return nil
		}
		return err
	}
	return nil
}

func (p *txPool) Send(tx models.GenericTransaction) {
	p.incomingTxsChan <- tx
}

func (p *txPool) Mempool() *Mempool {
	return p.mempool
}

func (p *txPool) RemoveFailedTxs(txErrors []models.TxError) error {
	if len(txErrors) == 0 {
		return nil
	}
	err := p.storage.SetTransactionErrors(txErrors...)
	if err != nil {
		return err
	}

	p.mempool.RemoveFailedTxs(txErrors)
	return nil
}

func earlierTimestamp(left, right *models.Timestamp) bool {
	if left == nil {
		return false
	}
	if right == nil {
		return true
	}
	return left.Before(*right)
}

func getReplacementError(txHash *common.Hash) models.TxError {
	return models.TxError{
		TxHash:       *txHash,
		ErrorMessage: ErrTxReplacementFailed.Error(),
	}
}

func getNonceError(txHash *common.Hash) models.TxError {
	return models.TxError{
		TxHash:       *txHash,
		ErrorMessage: ErrTxNonceTooLow.Error(),
	}
}
