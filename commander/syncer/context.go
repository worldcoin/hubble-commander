package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type SyncContext interface {
	SyncBatch(batch eth.DecodedBatch) error
	Commit() error
	Rollback(cause *error)
}

func NewSyncContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) SyncContext {
	switch batchType {
	case batchtype.Transfer, batchtype.Create2Transfer:
		return NewContext(storage, client, cfg, batchType)
	case batchtype.Deposit:
		return NewDepositsContext(storage, client, cfg)
	case batchtype.Genesis, batchtype.MassMigration:
		panic("invalid batch type")
	}
	return nil
}

type Context struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	tx        *db.TxController
	client    *eth.Client
	Syncer    TransactionSyncer
	BatchType batchtype.BatchType
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

func NewTestContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType) *Context {
	return newContext(storage, nil, client, cfg, batchType)
}

func newContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *Context {
	return &Context{
		cfg:       cfg,
		storage:   storage,
		tx:        tx,
		client:    client,
		Syncer:    NewTransactionSyncer(storage, client, batchType),
		BatchType: batchType,
	}
}

func (c *Context) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *Context) Rollback(cause *error) {
	c.tx.Rollback(cause)
}
