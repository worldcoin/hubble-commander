package db

import (
	"os"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if os.Getenv("HUBBLE_VERSION") == "" {
		os.Setenv("HUBBLE_VERSION", "dev-0.1.0")
	}

	if os.Getenv("HUBBLE_PORT") == "" {
		os.Setenv("HUBBLE_PORT", "8080")
	}
	
	if os.Getenv("HUBBLE_DBNAME") == "" {
		os.Setenv("HUBBLE_DBNAME", "hubble_test")
	}

	if os.Getenv("HUBBLE_DBUSER") == "" {
		os.Setenv("HUBBLE_DBUSER", "hubble")
	}
	
	if os.Getenv("HUBBLE_DBPASSWORD") == "" {
		os.Setenv("HUBBLE_DBPASSWORD", "root")
	}
	
	os.Exit(m.Run())
}

func TestGetDB(t *testing.T) {
	cfg := config.GetConfig()
	db, err := GetTestDB(cfg)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Ping())
}

func TestMigrations(t *testing.T) {
	cfg := config.GetConfig()
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
