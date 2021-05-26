package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

// transactionExecutor executes transactions & syncs batches. Manages a database transaction.
type transactionExecutor struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
}

// newTransactionExecutor creates a transactionExecutor and starts a database transaction.
func newTransactionExecutor(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (*transactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &transactionExecutor{
		cfg:     cfg,
		storage: txStorage,
		tx:      tx,
		client:  client,
	}, nil
}

// newTestTransactionExecutor creates a transactionExecutor without a database transaction.
func newTestTransactionExecutor(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) *transactionExecutor {
	return &transactionExecutor{
		cfg:     cfg,
		storage: storage,
		tx:      nil,
		client:  client,
	}
}

func (t *transactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *transactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}
