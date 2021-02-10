package db

import (
	"os"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if os.Getenv("Version") == "" {
		os.Setenv("Version", "dev-0.1.0")
	}

	if os.Getenv("Port") == "" {
		os.Setenv("Port", "8080")
	}
	
	if os.Getenv("DBName") == "" {
		os.Setenv("DBName", "hubble_test")
	}

	if os.Getenv("DBUser") == "" {
		os.Setenv("DBUser", "hubble")
	}
	
	if os.Getenv("DBPassword") == "" {
		os.Setenv("DBPassword", "root")
	}
	
	os.Exit(m.Run())
}

func TestGetDB(t *testing.T) {
	cfg, err := config.GetConfig()
	assert.NoError(t, err)

	db, err := GetTestDB(cfg)
	assert.NoError(t, err)
	defer db.Close()

	assert.NoError(t, db.Ping())
}

func TestMigrations(t *testing.T) {
	cfg, err := config.GetConfig()
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
