package mempool

import (
	"context"
	"sync"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
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

func NewTxPool(storage *st.Storage) (*txPool, error) {
	pool, err := NewMempool(storage)
	if err != nil {
		return nil, err
	}
	return &txPool{
		storage:         storage,
		mempool:         pool,
		incomingTxs:     make([]models.GenericTransaction, 0),
		incomingTxsChan: make(chan models.GenericTransaction, 1024),
	}, nil
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
	prevTxHash, err := p.mempool.AddOrReplace(p.storage, tx)
	if errors.Is(err, ErrTxReplacementFailed) {
		return p.storage.SetTransactionError(getReplacementError(&tx.GetBase().Hash))
	}
	if err != nil {
		return err
	}

	if prevTxHash != nil {
		return p.storage.RemovePendingTransaction(prevTxHash)
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

func getReplacementError(txHash *common.Hash) models.TxError {
	return models.TxError{
		TxHash:       *txHash,
		ErrorMessage: ErrTxReplacementFailed.Error(),
	}
}
