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
	cfg         *config.RollupConfig
	Storage     *st.Storage
	initStorage *st.Storage
	stateTree   *st.StateTree
	tx          *db.TxController
	client      *eth.Client
	ctx         context.Context
}

// NewTransactionExecutor creates a TransactionExecutor and starts a database transaction.
func NewTransactionExecutor(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	ctx context.Context,
) (*TransactionExecutor, error) {
	tx, txStorage, err := storage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return nil, err
	}

	return &TransactionExecutor{
		cfg:         cfg,
		Storage:     txStorage,
		initStorage: storage,
		stateTree:   st.NewStateTree(txStorage),
		tx:          tx,
		client:      client,
		ctx:         ctx,
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
		cfg:         cfg,
		Storage:     storage,
		initStorage: storage,
		stateTree:   st.NewStateTree(storage),
		tx:          nil,
		client:      client,
		ctx:         ctx,
	}
}

func (t *TransactionExecutor) Commit() error {
	return t.tx.Commit()
}

// nolint:gocritic
func (t *TransactionExecutor) Rollback(cause *error) {
	t.tx.Rollback(cause)
}

func (t *TransactionExecutor) RestartTransaction() error {
	tx, txStorage, err := t.initStorage.BeginTransaction(st.TxOptions{Postgres: true, Badger: true})
	if err != nil {
		return err
	}
	t.Storage = txStorage
	t.tx = tx
	t.stateTree = st.NewStateTree(txStorage)
	return nil
}
