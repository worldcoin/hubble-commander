package mempool

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
)

type testTxPool struct{}

func NewTestTxPool() TxPool {
	return &testTxPool{}
}

func (p *testTxPool) Send(models.GenericTransaction) {}

func (p *testTxPool) ReadTxs(context.Context) error {
	return nil
}

func (p *testTxPool) UpdateMempool() error {
	return nil
}

func (p *testTxPool) Mempool() *Mempool {
	return nil
}
