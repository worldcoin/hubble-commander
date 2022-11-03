package storage

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/utils"
)

type Storage struct {
	*BatchStorage
	*CommitmentStorage
	*ChainStateStorage
	*TransactionStorage
	*DepositStorage
	*RegisteredTokenStorage
	*RegisteredSpokeStorage
	*PendingStakeWithdrawalStorage
	StateTree           *StateTree
	AccountTree         *AccountTree
	database            *Database
	feeReceiverStateIDs map[string]uint32 // token ID => state id
}

type TxOptions struct {
	ReadOnly bool
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	database, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	return newStorageFromDatabase(database)
}

func newStorageFromDatabase(database *Database) (*Storage, error) {
	batchStorage, err := NewBatchStorage(database)
	if err != nil {
		return nil, err
	}

	commitmentStorage := NewCommitmentStorage(database)

	transactionStorage := NewTransactionStorage(database)

	depositStorage := NewDepositStorage(database)

	chainStateStorage := NewChainStateStorage(database)

	registeredTokenStorage := NewRegisteredTokenStorage(database)

	registeredSpokeStorage := NewRegisteredSpokeStorage(database)

	accountTree, err := NewAccountTree(database)
	if err != nil {
		return nil, err
	}

	pendingStakeWithdrawalStorage := NewPendingStakeWithdrawalStorage(database)

	storage := &Storage{
		BatchStorage:                  batchStorage,
		CommitmentStorage:             commitmentStorage,
		TransactionStorage:            transactionStorage,
		DepositStorage:                depositStorage,
		ChainStateStorage:             chainStateStorage,
		RegisteredTokenStorage:        registeredTokenStorage,
		RegisteredSpokeStorage:        registeredSpokeStorage,
		StateTree:                     NewStateTree(database),
		AccountTree:                   accountTree,
		PendingStakeWithdrawalStorage: pendingStakeWithdrawalStorage,
		database:                      database,
		feeReceiverStateIDs:           make(map[string]uint32),
	}
	err = storage.initBatchedTxsCounter()
	if err != nil {
		return nil, err
	}
	return storage, nil
}

func (s *Storage) copyWithNewDatabase(database *Database) *Storage {
	batchStorage := s.BatchStorage.copyWithNewDatabase(database)

	commitmentStorage := s.CommitmentStorage.copyWithNewDatabase(database)

	transactionStorage := s.TransactionStorage.copyWithNewDatabase(database)

	depositStorage := s.DepositStorage.copyWithNewDatabase(database)

	chainStateStorage := s.ChainStateStorage.copyWithNewDatabase(database)

	registeredTokenStorage := s.RegisteredTokenStorage.copyWithNewDatabase(database)

	registeredSpokeStorage := s.RegisteredSpokeStorage.copyWithNewDatabase(database)

	pendingStakeWithdrawalStorage := s.PendingStakeWithdrawalStorage.copyWithNewDatabase(database)

	stateTree := s.StateTree.copyWithNewDatabase(database)

	accountTree := s.AccountTree.copyWithNewDatabase(database)

	return &Storage{
		BatchStorage:                  batchStorage,
		CommitmentStorage:             commitmentStorage,
		TransactionStorage:            transactionStorage,
		DepositStorage:                depositStorage,
		ChainStateStorage:             chainStateStorage,
		RegisteredTokenStorage:        registeredTokenStorage,
		RegisteredSpokeStorage:        registeredSpokeStorage,
		PendingStakeWithdrawalStorage: pendingStakeWithdrawalStorage,
		StateTree:                     stateTree,
		AccountTree:                   accountTree,
		database:                      database,
		feeReceiverStateIDs:           utils.CopyStringUint32Map(s.feeReceiverStateIDs),
	}
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage) {
	txController, txDatabase := s.database.BeginTransaction(opts)

	return txController, s.copyWithNewDatabase(txDatabase)
}

func (s *Storage) ExecuteInReadWriteTransaction(fn func(txStorage *Storage) error) error {
	opts := TxOptions{ReadOnly: false}
	return s.ExecuteInTransaction(opts, fn)
}

func (s *Storage) ExecuteInTransaction(opts TxOptions, fn func(txStorage *Storage) error) error {
	return s.database.ExecuteInTransaction(opts, func(txDatabase *Database) error {
		return fn(s.copyWithNewDatabase(txDatabase))
	})
}

func (s *Storage) ExecuteInReadWriteTransactionWithSpan(
	ctx context.Context,
	fn func(txCtx context.Context, txStorage *Storage) error,
) error {
	opts := TxOptions{ReadOnly: false}
	return s.ExecuteInTransactionWithSpan(ctx, opts, fn)
}

func (s *Storage) ExecuteInTransactionWithSpan(
	ctx context.Context,
	opts TxOptions,
	fn func(txCtx context.Context, txStorage *Storage) error,
) error {
	return s.database.ExecuteInTransactionWithSpan(ctx, opts, func(txCtx context.Context, txDatabase *Database) error {
		return fn(txCtx, s.copyWithNewDatabase(txDatabase))
	})
}

func (s *Storage) TriggerGC() error {
	return s.database.Badger.TriggerGC()
}

func (s *Storage) Close() error {
	return s.database.Close()
}
