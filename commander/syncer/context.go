package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

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

	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return newContext(txStorage, tx, client, newTxsContext(txStorage, client, cfg, batchType))
	case batchtype.Deposit:
		return newContext(txStorage, tx, client, newDepositsContext(txStorage, client, cfg))
	case batchtype.Genesis, batchtype.MassMigration:
		panic("invalid batch type")
	}
	return nil
}

func NewTestContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	syncCtx batchContext,
) *Context {
	return newContext(storage, tx, client, syncCtx)
}

func (c *Context) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *Context) Rollback(cause *error) {
	c.tx.Rollback(cause)
}

func newContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	syncCtx batchContext,
) *Context {
	return &Context{
		storage:  storage,
		tx:       tx,
		client:   client,
		batchCtx: syncCtx,
	}
}
