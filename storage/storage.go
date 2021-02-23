package storage

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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

func (storage *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := sq.
		Insert("transaction").
		Values(
			tx.Hash,
			tx.FromIndex,
			tx.ToIndex,
			tx.Amount,
			tx.Fee,
			tx.Nonce,
			tx.Signature,
		).
		RunWith(storage.DB).
		PlaceholderFormat(sq.Dollar).
		Exec()

	return err
}

func (storage *Storage) GetTransaction(hash common.Hash) (*models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := storage.DB.Select(&res, "SELECT * FROM transaction WHERE tx_hash = $1", hash)
	if err != nil {
		return nil, err
	}
	return &res[0], nil
}
