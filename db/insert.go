package db

import (
	"database/sql"
)

type SqlBuilder interface {
	ToSql() (string, []interface{}, error)
}

func (d *Database) ExecBuilder(query SqlBuilder) (sql.Result, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return d.Exec(sqlQuery, args...)
}
