package db

import (
	"database/sql"
)

type SQLBuilder interface {
	ToSql() (string, []interface{}, error)
}

func (d *Database) ExecBuilder(query SQLBuilder) (sql.Result, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return d.Exec(sqlQuery, args...)
}
