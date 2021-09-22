package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type SyncContext struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	tx        *db.TxController
	client    *eth.Client
	Syncer    TransactionSyncer
	BatchType batchtype.BatchType
}

func NewSyncContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) (*SyncContext, error) {
	//TODO-div: extract and reuse new type
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{})
	if err != nil {
		return nil, err
	}

	return newSyncContext(txStorage, tx, client, cfg, batchType), nil
}

func NewTestSyncContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType) *SyncContext {
	return newSyncContext(storage, nil, client, cfg, batchType)
}

func newSyncContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *SyncContext {
	return &SyncContext{
		cfg:       cfg,
		storage:   storage,
		tx:        tx,
		client:    client,
		Syncer:    NewTransactionSyncer(storage, client, batchType),
		BatchType: batchType,
	}
}

func (c *SyncContext) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *SyncContext) Rollback(cause *error) {
	c.tx.Rollback(cause)
}
