package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type SyncContext interface {
	SyncNewBatch(batch eth.DecodedBatch) error
	UpdateExistingBatch(batch eth.DecodedBatch) error
}

type Context struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	client    *eth.Client
	Syncer    TransactionSyncer
	BatchType batchtype.BatchType
}

func NewTestContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType) *Context {
	return newContext(storage, client, cfg, batchType)
}

func newContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *Context {
	return &Context{
		cfg:       cfg,
		storage:   storage,
		client:    client,
		Syncer:    NewTransactionSyncer(storage, client, batchType),
		BatchType: batchType,
	}
}
