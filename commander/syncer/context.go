package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type batchContext interface {
	SyncCommitments(batch eth.DecodedBatch) error
	UpdateExistingBatch(batch eth.DecodedBatch) error
}

type Context struct {
	storage  *st.Storage
	tx       *db.TxController
	client   *eth.Client
	batchCtx batchContext
}

func NewContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *Context {
	tx, txStorage := storage.BeginTransaction(st.TxOptions{})
	return newContext(txStorage, tx, client, cfg, batchType)
}

func NewTestContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *Context {
	return newContext(storage, nil, client, cfg, batchType)
}

func newContext(
	txStorage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *Context {
	var batchCtx batchContext
	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		batchCtx = newTxsContext(txStorage, client, cfg, batchType)
	case batchtype.Deposit:
		batchCtx = newDepositsContext(txStorage, client)
	case batchtype.Genesis, batchtype.MassMigration:
		panic("invalid batch type")
	}
	return &Context{
		storage:  txStorage,
		tx:       tx,
		client:   client,
		batchCtx: batchCtx,
	}
}

func (c *Context) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *Context) Rollback(cause *error) {
	c.tx.Rollback(cause)
}
