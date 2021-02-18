package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func (storage *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := sq.
		Insert("transaction").
		Values(
			tx.Hash.String(),
			tx.FromIndex.Text(10),
			tx.ToIndex.Text(10),
			tx.Amount.Text(10),
			tx.Fee.Text(10),
			tx.Nonce.Text(10),
			tx.Signature,
		).
		RunWith(storage.db).
		PlaceholderFormat(sq.Dollar).
		Exec()

	return err
}
