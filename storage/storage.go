package storage

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB *sqlx.DB
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	dbInstance, err := db.GetDB(cfg)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: dbInstance}, nil
}
