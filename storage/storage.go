package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
)

type Storage struct {
	*StorageBase
	StateTree   *StateTree
	AccountTree *AccountTree
}

type StorageBase struct {
	Database            *Database
	feeReceiverStateIDs map[string]uint32 // token ID => state id
	latestBlockNumber   uint32
	syncedBlock         *uint64
}

type TxOptions struct {
	Postgres bool
	Badger   bool
	ReadOnly bool
}

func NewConfiguredStorage(cfg *config.Config) (storage *Storage, err error) {
	database, err := NewDatabase(cfg)
	if err != nil {
		return nil, err
	}

	storageBase := &StorageBase{
		Database:            database,
		feeReceiverStateIDs: make(map[string]uint32),
	}

	return &Storage{
		StorageBase: storageBase,
		StateTree:   NewStateTree(storageBase),
		AccountTree: NewAccountTree(storageBase),
	}, nil
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
	txController, txDatabase, err := s.Database.beginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txStorageBase := *s.StorageBase
	txStorageBase.Database = txDatabase

	txStorage := &Storage{
		StorageBase: &txStorageBase,
		StateTree:   NewStateTree(&txStorageBase),
		AccountTree: NewAccountTree(&txStorageBase),
	}

	return txController, txStorage, nil
}

// TODO-DB move to different file
// TODO-DB figure out how to make this solution more sustainable. Right now because of this
// 		   there are two BeginTransaction methods available on Storage
func (s *StorageBase) BeginTransaction(opts TxOptions) (*db.TxController, *StorageBase, error) {
	txController, txDatabase, err := s.Database.beginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txStorageBase := *s
	txStorageBase.Database = txDatabase

	return txController, &txStorageBase, nil
}

func (s *Storage) Close() error {
	return s.Database.Close()
}
