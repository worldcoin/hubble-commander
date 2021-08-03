package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
)

type Storage struct {
	*StorageBase
	*BatchStorage
	database    *Database
	StateTree   *StateTree
	AccountTree *AccountTree
}

type StorageBase struct {
	database            *Database
	feeReceiverStateIDs map[string]uint32 // token ID => state id
	latestBlockNumber   uint32
	syncedBlock         *uint64
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

	storageBase := &StorageBase{
		database:            database,
		feeReceiverStateIDs: make(map[string]uint32),
	}

	batchStorage := &BatchStorage{
		database: database,
	}

	return &Storage{
		StorageBase:  storageBase,
		BatchStorage: batchStorage,
		database:     database,
		StateTree:    NewStateTree(database),
		AccountTree:  NewAccountTree(database),
	}, nil
}

func (s *StorageBase) beginStorageBaseTransaction(opts TxOptions) (*db.TxController, *StorageBase, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txStorageBase := *s
	txStorageBase.database = txDatabase

	return txController, &txStorageBase, nil
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
	txController, txStorageBase, err := s.StorageBase.beginStorageBaseTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txStorage := &Storage{
		StorageBase: txStorageBase,
		StateTree:   NewStateTree(txStorageBase.database),
		AccountTree: NewAccountTree(txStorageBase.database),
	}

	return txController, txStorage, nil
}

func (s *Storage) Close() error {
	return s.database.Close()
}
