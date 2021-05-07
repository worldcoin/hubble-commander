package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
)

type QueryBuilder struct {
	db   *Database
	sql  string
	args []interface{}
	err  error
}

type SQLBuilder interface {
	ToSql() (string, []interface{}, error)
}

func (qb *QueryBuilder) Into(dest interface{}) error {
	if qb.err != nil {
		return qb.err
	}
	err := qb.db.Select(dest, qb.sql, qb.args...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (qb *QueryBuilder) Exec() (sql.Result, error) {
	if qb.err != nil {
		return nil, qb.err
	}
	sqlRes, err := qb.db.Exec(qb.sql, qb.args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sqlRes, nil
}

func (d *Database) Query(query SQLBuilder) *QueryBuilder {
	sqlQuery, args, err := query.ToSql()
	return &QueryBuilder{d, sqlQuery, args, errors.WithStack(err)}
}
