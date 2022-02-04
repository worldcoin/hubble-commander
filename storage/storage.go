package storage

import (
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

	return &Storage{
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
	}, nil
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

func (s *Storage) ExecuteInTransaction(opts TxOptions, fn func(txStorage *Storage) error) error {
	return s.database.ExecuteInTransaction(opts, func(txDatabase *Database) error {
		return fn(s.copyWithNewDatabase(txDatabase))
	})
}

func (s *Storage) Close() error {
	return s.database.Close()
}
