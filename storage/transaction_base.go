package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

type TransactionStorage struct {
	database *Database
}

func NewTransactionStorage(database *Database) *TransactionStorage {
	return &TransactionStorage{
		database: database,
	}
}

func (s *TransactionStorage) copyWithNewDatabase(database *Database) *TransactionStorage {
	newTransactionStorage := *s
	newTransactionStorage.database = database

	return &newTransactionStorage
}

func (s *TransactionStorage) BeginTransaction(opts TxOptions) (*db.TxController, *TransactionStorage, error) {
	txController, txDatabase, err := s.database.BeginTransaction(opts)
	if err != nil {
		return nil, nil, err
	}

	txTransactionStorage := *s
	txTransactionStorage.database = txDatabase

	return txController, &txTransactionStorage, nil
}

func (s *TransactionStorage) GetLatestTransactionNonce(accountStateID uint32) (*models.Uint256, error) {
	var tx models.StoredTransaction
	err := s.database.Badger.Iterator(models.StoredTransactionPrefix, badger.ReversePrefetchIteratorOpts,
		func(item *bdg.Item) (bool, error) {
			err := item.Value(tx.SetBytes)
			if err != nil {
				return false, err
			}
			return tx.FromStateID == accountStateID, nil
		})
	if err == badger.ErrIteratorFinished {
		return nil, NewNotFoundError("transaction")
	}
	if err != nil {
		return nil, err
	}
	return &tx.Nonce, nil
}

// BatchMarkTransactionAsIncluded TODO-tx: replace usage with custom functions
func (s *TransactionStorage) BatchMarkTransactionAsIncluded(txHashes []common.Hash, commitmentID *models.CommitmentID) error {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txHashes {
		var storedTx models.StoredTransaction
		err = txStorage.database.Badger.Get(txHashes[i], &storedTx)
		if err == bh.ErrNotFound {
			return NewNotFoundError("transaction")
		}
		if err != nil {
			return err
		}

		storedTx.CommitmentID = commitmentID
		err = txStorage.database.Badger.Update(storedTx.Hash, storedTx)
		if err == bh.ErrNotFound {
			return NewNotFoundError("transaction")
		}
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *TransactionStorage) SetTransactionError(txHash common.Hash, errorMessage string) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("transaction_base").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("error_message", errorMessage),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) GetTransactionCount() (*int, error) {
	latestCommitment, err := s.GetLatestCommitment()
	if IsNotFoundError(err) {
		return ref.Int(0), nil
	}
	if err != nil {
		return nil, err
	}

	res := make([]int, 0, 1)
	err = s.database.Postgres.Query(
		s.database.QB.Select("COUNT(1)").
			From("transaction_base").
			Where(squirrel.LtOrEq{"batch_id": latestCommitment.ID.BatchID}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return ref.Int(0), nil
	}
	return &res[0], nil
}

func (s *TransactionStorage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	res := make([]common.Hash, 0, 32*len(batchIDs))
	err := s.database.Postgres.Query(
		s.database.QB.Select("transaction_base.tx_hash").
			From("transaction_base").
			Where(squirrel.Eq{"batch_id": batchIDs}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return res, nil
}
