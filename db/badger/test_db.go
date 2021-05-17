package badger

import (
	"github.com/dgraph-io/badger/v3"
	bh "github.com/timshannon/badgerhold/v3"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	options := bh.DefaultOptions
	options.Options = badger.DefaultOptions("").
		WithInMemory(true).
		WithLoggingLevel(badger.WARNING)

	store, err := bh.Open(options)
	if err != nil {
		return nil, err
	}
	db := &Database{store: store}
	teardown := func() error {
		return db.Close()
	}

	return &TestDB{
		DB:       db,
		Teardown: teardown,
	}, nil
}
