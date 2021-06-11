package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

// transactionExecutor executes transactions & syncs batches. Manages a database transaction.
type transactionExecutor struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	stateTree *st.StateTree
	tx        *db.TxController
	client    *eth.Client
	ctx       context.Context
}

// newTransactionExecutor creates a transactionExecutor and starts a database transaction.
func newTransactionExecutor(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) (*transactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &transactionExecutor{
		cfg:       cfg,
		storage:   txStorage,
		stateTree: st.NewStateTree(txStorage),
		tx:        tx,
		client:    client,
	}, nil
}

// newTransactionExecutor creates a transactionExecutor with context and starts a database transaction.
func newTransactionExecutorWithCtx(
	ctx context.Context,
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
) (*transactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &transactionExecutor{
		cfg:       cfg,
		storage:   txStorage,
		stateTree: st.NewStateTree(txStorage),
		tx:        tx,
		client:    client,
		ctx:       ctx,
	}, nil
}

// newTestTransactionExecutor creates a transactionExecutor without a database transaction.
func newTestTransactionExecutor(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) *transactionExecutor {
	return &transactionExecutor{
		cfg:       cfg,
		storage:   storage,
		stateTree: st.NewStateTree(storage),
		tx:        nil,
		client:    client,
		ctx:       context.Background(),
	}
}

func (t *transactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *transactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}
