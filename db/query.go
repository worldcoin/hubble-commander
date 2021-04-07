package db

import (
	"database/sql"
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

func (qb *QueryBuilder) Exec() (sql.Result, error) {
	if qb.err != nil {
		return nil, qb.err
	}
	sqlRes, err := qb.db.Exec(qb.sql, qb.args...)
	if err != nil {
		return nil, err
	}
	return sqlRes, nil
}

func (d *Database) Query(query SQLBuilder) *QueryBuilder {
	sql, args, err := query.ToSql()
	return &QueryBuilder{d, sql, args, err}
}
