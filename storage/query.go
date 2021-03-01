package storage

import "github.com/Masterminds/squirrel"

type QueryBuilder struct {
	storage *Storage
	sql     string
	args    []interface{}
	err     error
}

func (qb *QueryBuilder) Into(dest interface{}) error {
	if qb.err != nil {
		return qb.err
	}
	err := qb.storage.DB.Select(dest, qb.sql, qb.args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Query(query squirrel.SelectBuilder) *QueryBuilder {
	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	return &QueryBuilder{s, sql, args, err}
}
