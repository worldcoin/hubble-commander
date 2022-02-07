package executor

import (
	"context"

	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	st "github.com/Worldcoin/hubble-commander/storage"
)

type ExecutionContext struct {
	cfg              *config.RollupConfig
	storage          *st.Storage
	tx               *db.TxController
	client           *eth.Client
	txsSender        TxsSender
	ctx              context.Context
	commanderMetrics *metrics.CommanderMetrics
	Applier          *applier.Applier
}

type TxsSender interface {
	TransferTxsSender
	C2TTxsSender
	MassMigrationTxsSender
}

// NewExecutionContext creates a ExecutionContext and starts a database transaction.
func NewExecutionContext(
	storage *st.Storage,
	client *eth.Client,
	txsSender TxsSender,
	cfg *config.RollupConfig,
	commanderMetrics *metrics.CommanderMetrics,
	ctx context.Context,
) *ExecutionContext {
	tx, txStorage := storage.BeginTransaction(st.TxOptions{})

	return &ExecutionContext{
		cfg:              cfg,
		storage:          txStorage,
		tx:               tx,
		txsSender:        txsSender,
		client:           client,
		ctx:              ctx,
		commanderMetrics: commanderMetrics,
		Applier:          applier.NewApplier(txStorage),
	}
}

// NewTestExecutionContext creates a ExecutionContext without a database transaction.
func NewTestExecutionContext(storage *st.Storage, client *eth.Client, txsSender TxsSender, cfg *config.RollupConfig) *ExecutionContext {
	return &ExecutionContext{
		cfg:              cfg,
		storage:          storage,
		tx:               nil,
		client:           client,
		txsSender:        txsSender,
		ctx:              context.Background(),
		commanderMetrics: metrics.NewCommanderMetrics(),
		Applier:          applier.NewApplier(storage),
	}
}

func (c *ExecutionContext) Commit() error {
	return c.tx.Commit()
}

// nolint:gocritic
func (c *ExecutionContext) Rollback(cause *error) {
	c.tx.Rollback(cause)
}
