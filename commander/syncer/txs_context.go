package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

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

func NewMempoolContext(pool *mempool.Mempool) *MempoolContext {
	txController, txMempool := pool.BeginTransaction()
	return &MempoolContext{
		Mempool:      txMempool,
		txController: txController,
	}
}

type TxsContext struct {
	cfg        *config.RollupConfig
	storage    *st.Storage
	client     *eth.Client
	mempoolCtx *MempoolContext
	Syncer     TransactionSyncer
	TxType     txtype.TransactionType
}

func NewTestTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	txType txtype.TransactionType,
) (*TxsContext, error) {
	pool, err := mempool.NewMempool(storage)
	if err != nil {
		return nil, err
	}
	return newTxsContext(storage, client, pool, cfg, txType), nil
}

func newTxsContext(
	storage *st.Storage,
	client *eth.Client,
	pool *mempool.Mempool,
	cfg *config.RollupConfig,
	txType txtype.TransactionType,
) *TxsContext {
	return &TxsContext{
		cfg:        cfg,
		storage:    storage,
		client:     client,
		mempoolCtx: NewMempoolContext(pool),
		Syncer:     NewTransactionSyncer(storage, txType),
		TxType:     txType,
	}
}
