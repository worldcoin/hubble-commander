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
	sql, args, err := sq.
		Insert("users").Columns("name", "age").
		Values("moe", 13).Values("larry", sq.Expr("? + 5", 12)).
		ToSql()
}
