package syncer

import "github.com/Worldcoin/hubble-commander/mempool"

type MempoolContext struct {
	Mempool      *mempool.TxMempool
	txController *mempool.TxController
}

func NewMempoolContext(pool *mempool.Mempool) *MempoolContext {
	txController, txMempool := pool.BeginTransaction()
	return &MempoolContext{
		Mempool:      txMempool,
		txController: txController,
	}
}

func (c *MempoolContext) Commit() {
	c.txController.Commit()
}

func (c *MempoolContext) Rollback() {
	c.txController.Rollback()
}
