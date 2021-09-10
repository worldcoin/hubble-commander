package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
)

// ExecutionContext executes transactions & syncs batches. Manages a database transaction.
type ExecutionContext struct {
	cfg     *config.RollupConfig
	storage *st.Storage
	tx      *db.TxController
	client  *eth.Client
	ctx     context.Context
}

// NewExecutionContext creates a ExecutionContext and starts a database transaction.
func NewExecutionContext(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) (*ExecutionContext, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{})
	if err != nil {
		return nil, err
	}

	return &ExecutionContext{
		cfg:     cfg,
		storage: txStorage,
		tx:      tx,
		client:  client,
		ctx:     ctx,
	}, nil
}

// NewTestExecutionContext creates a ExecutionContext without a database transaction.
func NewTestExecutionContext(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) *ExecutionContext {
	return &ExecutionContext{
		cfg:     cfg,
		storage: storage,
		tx:      nil,
		client:  client,
		ctx:     context.Background(),
	}
}

func (t *ExecutionContext) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *ExecutionContext) Rollback(cause *error) {
	t.tx.Rollback(cause)
}
