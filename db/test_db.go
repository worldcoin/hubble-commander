package db

import (
	"bytes"

	"github.com/pkg/errors"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	db, err := NewInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		return db.Close()
	}

	return &TestDB{
		DB:       db,
		Teardown: teardown,
	}, nil
}

func (d *TestDB) Clone() (*TestDB, error) {
	clonedBadger, err := NewTestDB()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var backup bytes.Buffer
	_, err = d.DB.store.Badger().Backup(&backup, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = clonedBadger.DB.store.Badger().Load(&backup, 16)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return clonedBadger, nil
}
