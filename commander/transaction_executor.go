package commander

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type transactionExecutorOpts struct {
	ctx context.Context
	// Ignore nonce field on transfers, assign it with the correct value from the state.
	AssumeNonces bool
}

// transactionExecutor executes transactions & syncs batches. Manages a database transaction.
type transactionExecutor struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	stateTree *st.StateTree
	tx        *db.TxController
	client    *eth.Client
	opts      transactionExecutorOpts
}

// newTransactionExecutor creates a transactionExecutor and starts a database transaction.
func newTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	opts transactionExecutorOpts,
) (*transactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	defaultOpts(&opts)
	return &transactionExecutor{
		cfg:       cfg,
		storage:   txStorage,
		stateTree: st.NewStateTree(txStorage),
		tx:        tx,
		client:    client,
		opts:      opts,
	}, nil
}

// newTestTransactionExecutor creates a transactionExecutor without a database transaction.
func newTestTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	opts transactionExecutorOpts,
) *transactionExecutor {
	defaultOpts(&opts)
	return &transactionExecutor{
		cfg:       cfg,
		storage:   storage,
		stateTree: st.NewStateTree(storage),
		tx:        nil,
		client:    client,
		opts:      opts,
	}
}

func (t *transactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *transactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}

func defaultOpts(opts *transactionExecutorOpts) {
	if opts.ctx == nil {
		opts.ctx = context.Background()
	}
}
