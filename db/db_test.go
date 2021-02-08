package db

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDB(t *testing.T) {
	cfg, err := config.GetConfig("../config.template.yaml")
	if err != nil {
		t.Fatal(err)
	}

	db, err := GetDB(cfg)
	if err != nil {
		t.Fatal(err)
	}
	
	assert.NoError(t, db.Ping())
}
