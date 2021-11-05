package syncer

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type DepositsContext struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
}

func NewDepositsContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) *DepositsContext {
	tx, txStorage := storage.BeginTransaction(st.TxOptions{})
	return newDepositsContext(txStorage, tx, client, cfg)
}

func NewTestDepositsContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) *DepositsContext {
	return newDepositsContext(storage, nil, client, cfg)
}

func newDepositsContext(
	storage *st.Storage,
	tx *db.TxController,
	client *eth.Client,
	cfg *config.RollupConfig,
) *DepositsContext {
	return &DepositsContext{
		cfg:     cfg,
		storage: storage,
		tx:      tx,
		client:  client,
	}
}

func (c *DepositsContext) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *DepositsContext) Rollback(cause *error) {
	c.tx.Rollback(cause)
}
