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
	sql, args, err := s.QB.Select("*").
		From("transaction").
		Where(squirrel.Eq{"tx_hash": hash}).
		ToSql()
	if err != nil {
		return nil, err
	}
	err = s.DB.Select(&res, sql, args...)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}
