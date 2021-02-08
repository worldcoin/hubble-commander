package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDB(t *testing.T) {
	cfg, err := config.GetConfig("../config.template.yaml")
	assert.NoError(t, err)

	db, err := GetTestDB(cfg)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Ping())
}

func TestMigrations(t *testing.T) {
	cfg, err := config.GetConfig("../config.template.yaml")
	assert.NoError(t, err)

	db, err := GetTestDB(cfg)
	assert.NoError(t, err)
	defer db.Close()

	migrator, err := GetMigrator(cfg)
	assert.NoError(t, err)

	assert.NoError(t, migrator.Up())

	_, err = sq.Select("*").From("person").
		RunWith(db).Query()
	assert.NoError(t, err)

	assert.NoError(t, migrator.Down())

	_, err = sq.Select("*").From("person").
		RunWith(db).Query()
	assert.Error(t, err)

}
