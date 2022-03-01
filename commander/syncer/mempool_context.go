package syncer

import "github.com/Worldcoin/hubble-commander/mempool"

type MempoolContext struct {
	Mempool      *mempool.TxMempool
	txController *mempool.TxController
}

func (c *MempoolContext) Commit() {
	c.txController.Commit()
}

func (c *MempoolContext) Rollback() {
	c.txController.Rollback()
}
