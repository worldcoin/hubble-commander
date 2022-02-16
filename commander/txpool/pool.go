package txpool

import (
	"context"

	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TxPool struct {
	storage *st.Storage
	Pool    *mempool.Mempool
	TxChan  chan models.GenericTransaction
}

func NewTxPool(storage *st.Storage) (*TxPool, error) {
	pool, err := mempool.NewMempool(storage)
	if err != nil {
		return nil, err
	}
	return &TxPool{
		storage: storage,
		Pool:    pool,
		TxChan:  make(chan models.GenericTransaction, 1024),
	}, nil
}

func (p *TxPool) ReadTxs(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():

		case tx := <-p.TxChan:
			//TODO: introduce nonce cache or pass storage to AddOrReplace
			stateLeaf, err := p.storage.StateTree.Leaf(tx.GetFromStateID())
			if err != nil {
				return err
			}
			err = p.Pool.AddOrReplace(tx, stateLeaf.Nonce.Uint64())
			if err == mempool.ErrTxReplacementFailed {
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
