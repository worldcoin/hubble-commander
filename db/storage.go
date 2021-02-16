package db

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func (storage *Storage) AddTransaction(tx *models.Transaction) error {
	sql, args, err := sq.
		Insert("transaction").
		Values(
			tx.Hash.String(),
			tx.FromIndex,
			tx.ToIndex,
			tx.Amount,
			tx.Fee,
			tx.Nonce,
			tx.Signature,
		).ToSql()
	if err != nil {
		return err
	}
	fmt.Println(sql)
	fmt.Println(args)
	return nil
}
