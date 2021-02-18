package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestGetDB(t *testing.T) {
	cfg := config.GetTestConfig()
	db, err := GetTestDB(&cfg)
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()
	require.NoError(t, db.Ping())
}

func TestMigrations(t *testing.T) {
	cfg := config.GetTestConfig()
	db, err := GetTestDB(&cfg)
	require.NoError(t, err)
	defer func() { require.NoError(t, db.Close()) }()

	migrator, err := GetMigrator(&cfg)
	require.NoError(t, err)

	require.NoError(t, migrator.Up())

	_, err = sq.Select("*").From("transaction").
		RunWith(db).Query()
	require.NoError(t, err)

	require.NoError(t, migrator.Down())

	_, err = sq.Select("*").From("transaction").
		RunWith(db).Query()
	require.Error(t, err)
}
