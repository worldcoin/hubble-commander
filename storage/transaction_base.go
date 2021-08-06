package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
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

func (s *TransactionStorage) addTransactionBase(txBase *models.TransactionBase, txType txtype.TransactionType) (*models.Timestamp, error) {
	res := make([]models.Timestamp, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Insert("transaction_base").
			Values(
				txBase.Hash,
				txType,
				txBase.FromStateID,
				txBase.Amount,
				txBase.Fee,
				txBase.Nonce,
				txBase.Signature,
				txBase.IncludedInCommitment,
				txBase.ErrorMessage,
				"NOW()",
			).
			Suffix("RETURNING receive_time"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (s *TransactionStorage) BatchAddTransactionBase(txs []models.TransactionBase) error {
	query := s.database.QB.Insert("transaction_base")
	for i := range txs {
		query = query.Values(
			txs[i].Hash,
			txs[i].TxType,
			txs[i].FromStateID,
			txs[i].Amount,
			txs[i].Fee,
			txs[i].Nonce,
			txs[i].Signature,
			txs[i].IncludedInCommitment,
			txs[i].ErrorMessage,
		)
	}
	res, err := s.database.Postgres.Query(query).Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *TransactionStorage) GetLatestTransactionNonce(accountStateID uint32) (*models.Uint256, error) {
	res := make([]models.Uint256, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("transaction_base.nonce").
			From("transaction_base").
			Where(squirrel.Eq{"from_state_id": accountStateID}).
			OrderBy("nonce DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return &res[0], nil
}

func (s *TransactionStorage) BatchMarkTransactionAsIncluded(txHashes []common.Hash, commitmentID *int32) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("transaction_base").
			Where(squirrel.Eq{"tx_hash": txHashes}).
			Set("included_in_commitment", commitmentID),
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
	res := make([]int, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select("COUNT(1)").
			From("transaction_base").
			Join("commitment on commitment.commitment_id = transaction_base.included_in_commitment").
			Where(squirrel.NotEq{"included_in_batch": nil}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) < 1 {
		return ref.Int(0), nil
	}
	return &res[0], nil
}

func (s *Storage) GetTransactionHashesByBatchIDs(batchIDs ...models.Uint256) ([]common.Hash, error) {
	res := make([]common.Hash, 0, 32*len(batchIDs))
	err := s.database.Postgres.Query(
		s.database.QB.Select("transaction_base.tx_hash").
			From("transaction_base").
			Join("commitment on commitment.commitment_id = transaction_base.included_in_commitment").
			Where(squirrel.Eq{"commitment.included_in_batch": batchIDs}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return res, nil
}
