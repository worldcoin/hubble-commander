package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type TransactionExecutorOpts struct {
	Ctx context.Context
	// Ignore nonce field on transfers, assign it with the correct value from the state.
	AssumeNonces bool // TODO-AFS remove
}

// TransactionExecutor executes transactions & syncs batches. Manages a database transaction.
type TransactionExecutor struct {
	cfg       *config.RollupConfig
	storage   *st.Storage
	stateTree *st.StateTree
	tx        *db.TxController
	client    *eth.Client
	opts      TransactionExecutorOpts
}

// NewTransactionExecutor creates a TransactionExecutor and starts a database transaction.
func NewTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	opts TransactionExecutorOpts,
) (*TransactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	defaultOpts(&opts)
	return &TransactionExecutor{
		cfg:       cfg,
		storage:   txStorage,
		stateTree: st.NewStateTree(txStorage),
		tx:        tx,
		client:    client,
		opts:      opts,
	}, nil
}

// NewTestTransactionExecutor creates a TransactionExecutor without a database transaction.
func NewTestTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	opts TransactionExecutorOpts,
) *TransactionExecutor {
	defaultOpts(&opts)
	return &TransactionExecutor{
		cfg:       cfg,
		storage:   storage,
		stateTree: st.NewStateTree(storage),
		tx:        nil,
		client:    client,
		opts:      opts,
	}
}

func (t *TransactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *TransactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}

func defaultOpts(opts *TransactionExecutorOpts) {
	if opts.Ctx == nil {
		opts.Ctx = context.Background()
	}
}
