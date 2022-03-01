package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

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
	txPool, err := mempool.NewTxPool(storage)
	if err != nil {
		return nil, err
	}
	return newTxsContext(storage, client, txPool.Mempool(), cfg, txType), nil
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

func (c *TxsContext) Commit() {
	c.mempoolCtx.Commit()
}

func (c *TxsContext) Rollback() {
	c.mempoolCtx.Rollback()
}
