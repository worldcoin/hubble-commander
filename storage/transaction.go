package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := s.QB.Insert("transaction").
		Values(
			tx.Hash,
			tx.FromIndex,
			tx.ToIndex,
			tx.Amount,
			tx.Fee,
			tx.Nonce,
			tx.Signature,
		).
		RunWith(s.DB).
		Exec()

	return err
}

func (s *Storage) GetTransaction(hash common.Hash) (*models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := s.Query(
		squirrel.Select("*").
			From("transaction").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

type QueryBuilder struct {
	storage *Storage
	sql     string
	args    []interface{}
	err     error
}

func (qb *QueryBuilder) Into(dest interface{}) error {
	if qb.err != nil {
		return qb.err
	}
	err := qb.storage.DB.Select(dest, qb.sql, qb.args...)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) Query(query squirrel.SelectBuilder) *QueryBuilder {
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	return &QueryBuilder{storage, sql, args, err}
}
