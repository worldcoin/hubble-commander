package db

import (
	"github.com/Masterminds/squirrel"
)

type QueryBuilder struct {
	db   *Database
	sql  string
	args []interface{}
	err  error
}

func (qb *QueryBuilder) Into(dest interface{}) error {
	if qb.err != nil {
		return qb.err
	}
	err := qb.db.Select(dest, qb.sql, qb.args...)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) Query(query squirrel.SelectBuilder) *QueryBuilder {
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	return &QueryBuilder{d, sql, args, err}
}

func (d *Database) InsertQuery(query squirrel.InsertBuilder) *QueryBuilder {
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	return &QueryBuilder{d, sql, args, err}
}
