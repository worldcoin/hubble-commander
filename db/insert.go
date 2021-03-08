package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

func (d *Database) Insert(query squirrel.InsertBuilder) (sql.Result, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return d.Exec(sqlQuery, args...)
}
