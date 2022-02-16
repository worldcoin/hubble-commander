package mempool

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TxPool interface {
	Send(tx models.GenericTransaction)
	ReadTxs(ctx context.Context) error
	Mempool() *Mempool
}

type txPool struct {
	storage *st.Storage
	mempool *Mempool
	TxChan  chan models.GenericTransaction
}

func NewTxPool(storage *st.Storage) (TxPool, error) {
	pool, err := NewMempool(storage)
	if err != nil {
		return nil, err
	}
	return &txPool{
		storage: storage,
		mempool: pool,
		TxChan:  make(chan models.GenericTransaction, 1024),
	}, nil
}

func (p *txPool) ReadTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case tx := <-p.TxChan:
			//TODO: introduce nonce cache or pass storage to AddOrReplace
			stateLeaf, err := p.storage.StateTree.Leaf(tx.GetFromStateID())
			if err != nil {
				return err
			}
			err = p.mempool.AddOrReplace(tx, stateLeaf.Nonce.Uint64())
			if err == ErrTxReplacementFailed {
				continue
			}
			if err != nil {
				return err
			}
		default:
			// return if there is no more txs in p.TxChan
			return nil
		}
	}
}

func (p *txPool) Send(tx models.GenericTransaction) {
	p.TxChan <- tx
}

func (p *txPool) Mempool() *Mempool {
	return p.mempool
}
