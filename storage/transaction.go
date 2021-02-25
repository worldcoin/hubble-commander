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
	err := storage.Query(
		&res,
		squirrel.Select("*").
			From("transaction").
			Where(squirrel.Eq{"tx_hash": hash}),
	)
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}

func (storage *Storage) Query(dest interface{}, query squirrel.SelectBuilder) error {
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}
	err = storage.DB.Select(&dest, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
