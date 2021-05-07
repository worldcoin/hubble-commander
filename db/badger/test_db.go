package badger

import (
	"github.com/dgraph-io/badger/v3"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return nil, err
	}

	teardown := func() error {
		return db.Close()
	}

	return &TestDB{
		DB:       &Database{badger: db},
		Teardown: teardown,
	}, nil
}
