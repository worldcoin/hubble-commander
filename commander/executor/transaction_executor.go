package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

// TransactionExecutor executes transactions & syncs batches. Manages a database transaction.
type TransactionExecutor struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
	ctx     context.Context
}

// NewTransactionExecutor creates a TransactionExecutor and starts a database transaction.
func NewTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) (*TransactionExecutor, error) {
	tx, txStorageBase, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &TransactionExecutor{
		cfg: cfg,
		storage: &st.Storage{
			StorageBase: txStorageBase,
			StateTree:   st.NewStateTree(txStorageBase),
			AccountTree: st.NewAccountTree(txStorageBase),
		},
		tx:     tx,
		client: client,
		ctx:    ctx,
	}, nil
}

// NewTestTransactionExecutor creates a TransactionExecutor without a database transaction.
func NewTestTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) *TransactionExecutor {
	return &TransactionExecutor{
		cfg:     cfg,
		storage: storage,
		tx:      nil,
		client:  client,
		ctx:     ctx,
	}
}

func (t *TransactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *TransactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}
