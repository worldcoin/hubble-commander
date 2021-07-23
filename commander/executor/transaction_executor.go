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
	cfg       *config.RollupConfig
	storage   *st.InternalStorage
	stateTree *st.StateTree
	tx        *db.TxController
	client    *eth.Client
	ctx       context.Context
}

// TODO-INTERNAL revisit and change storage to st.Storage
// NewTransactionExecutor creates a TransactionExecutor and starts a database transaction.
func NewTransactionExecutor(
	storage *st.InternalStorage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) (*TransactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &TransactionExecutor{
		cfg:       cfg,
		storage:   txStorage,
		stateTree: st.NewStateTree(txStorage),
		tx:        tx,
		client:    client,
		ctx:       ctx,
	}, nil
}

// NewTestTransactionExecutor creates a TransactionExecutor without a database transaction.
func NewTestTransactionExecutor(
	storage *st.InternalStorage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) *TransactionExecutor {
	return &TransactionExecutor{
		cfg:       cfg,
		storage:   storage,
		stateTree: st.NewStateTree(storage),
		tx:        nil,
		client:    client,
		ctx:       ctx,
	}
}

func (t *TransactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *TransactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}
