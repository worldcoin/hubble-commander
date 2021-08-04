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

	batchStorage := &BatchStorage{
		database: database,
	}

	commitmentStorage := &CommitmentStorage{
		database: database,
	}

	transactionStorage := &TransactionStorage{
		database: database,
	}

	chainStateStorage := &ChainStateStorage{
		database: database,
	}

	return &Storage{
		BatchStorage:        batchStorage,
		CommitmentStorage:   commitmentStorage,
		TransactionStorage:  transactionStorage,
		ChainStateStorage:   chainStateStorage,
		StateTree:           NewStateTree(database),
		AccountTree:         NewAccountTree(database),
		database:            database,
		feeReceiverStateIDs: make(map[string]uint32),
	}, nil
}

func (s *Storage) BeginTransaction(opts TxOptions) (*db.TxController, *Storage, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txBatchStorage := *s.BatchStorage
	txBatchStorage.database = txDatabase

	txCommitmentStorage := *s.CommitmentStorage
	txCommitmentStorage.database = txDatabase

	txTransactionStorage := *s.TransactionStorage
	txTransactionStorage.database = txDatabase

	txChainStateStorage := *s.ChainStateStorage
	txChainStateStorage.database = txDatabase

	txStorage := &Storage{
		BatchStorage:        &txBatchStorage,
		CommitmentStorage:   &txCommitmentStorage,
		TransactionStorage:  &txTransactionStorage,
		ChainStateStorage:   &txChainStateStorage,
		StateTree:           NewStateTree(txDatabase),
		AccountTree:         NewAccountTree(txDatabase),
		database:            txDatabase,
		feeReceiverStateIDs: utils.CopyStringUint32Map(s.feeReceiverStateIDs),
	}

	return txController, txStorage, nil
}

func (s *Storage) Close() error {
	return s.database.Close()
}
