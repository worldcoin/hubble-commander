package badger

import (
	"bytes"

	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

type TestDB struct {
	DB       *Database
	Teardown func() error
}

func NewTestDB() (*TestDB, error) {
	options := bh.DefaultOptions
	options.Encoder = Encode
	options.Decoder = Decode
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
