package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TxsContext struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	client    *eth.Client
	Syncer    TransactionSyncer
	BatchType batchtype.BatchType
}

func NewTestTxsContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType) *TxsContext {
	return newTxsContext(storage, client, cfg, batchType)
}

func newTxsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) *TxsContext {
	return &TxsContext{
		cfg:       cfg,
		storage:   storage,
		client:    client,
		Syncer:    NewTransactionSyncer(storage, batchType),
		BatchType: batchType,
	}
}
