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
	*RegisteredTokenStorage
	StateTree           *StateTree
	AccountTree         *AccountTree
	database            *Database
	feeReceiverStateIDs map[string]uint32 // token ID => state id
}

type TxOptions struct {
	Postgres bool
	Badger   bool
	ReadOnly bool
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	database, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	storage := newStorageFromDatabase(database)

	return storage, nil
}

func newStorageFromDatabase(database *Database) *Storage {
	batchStorage := NewBatchStorage(database)

	commitmentStorage := NewCommitmentStorage(database)

	transactionStorage := NewTransactionStorage(database)

	chainStateStorage := NewChainStateStorage(database)

	registeredTokenStorage := NewRegisteredTokenStorage(database)

	return &Storage{
		BatchStorage:           batchStorage,
		CommitmentStorage:      commitmentStorage,
		TransactionStorage:     transactionStorage,
		ChainStateStorage:      chainStateStorage,
		RegisteredTokenStorage: registeredTokenStorage,
		StateTree:              NewStateTree(database),
		AccountTree:            NewAccountTree(database),
		database:               database,
		feeReceiverStateIDs:    make(map[string]uint32),
	}
}

func (s *Storage) copyWithNewDatabase(database *Database) *Storage {
	batchStorage := s.BatchStorage.copyWithNewDatabase(database)

	commitmentStorage := s.CommitmentStorage.copyWithNewDatabase(database)

	transactionStorage := s.TransactionStorage.copyWithNewDatabase(database)

	chainStateStorage := s.ChainStateStorage.copyWithNewDatabase(database)

	registeredTokenStorage := s.RegisteredTokenStorage.copyWithNewDatabase(database)

	stateTree := s.StateTree.copyWithNewDatabase(database)

	accountTree := s.AccountTree.copyWithNewDatabase(database)

	return &Storage{
		BatchStorage:           batchStorage,
		CommitmentStorage:      commitmentStorage,
		TransactionStorage:     transactionStorage,
		ChainStateStorage:      chainStateStorage,
		RegisteredTokenStorage: registeredTokenStorage,
		StateTree:              stateTree,
		AccountTree:            accountTree,
		database:               database,
		feeReceiverStateIDs:    utils.CopyStringUint32Map(s.feeReceiverStateIDs),
	}
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txStorage := s.copyWithNewDatabase(txDatabase)

	return txController, txStorage, nil
}

func (s *Storage) Close() error {
	return s.database.Close()
}
