package storage

import (
	"github.com/Worldcoin/hubble-commander/db/postgres"
)

type TestStorage struct {
	*Storage
	Teardown func() error
}

func NewTestStorage() (*TestStorage, error) {
	testDB, err := postgres.NewTestDB()
	if err != nil {
		return nil, err
	}

	return &TestStorage{
		Storage: &Storage{
			Postgres: testDB.DB,
			QB:       getQueryBuilder(),
		},
		Teardown: func() error {
			return testDB.Teardown()
		},
	}, nil
}
