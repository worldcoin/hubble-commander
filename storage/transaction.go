package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := s.DB.Insert(
		s.QB.Insert("transaction").
			Values(
				tx.Hash,
				tx.FromIndex,
				tx.ToIndex,
				tx.Amount,
				tx.Fee,
				tx.Nonce,
				tx.Signature,
			),
	)

	return err
}

func (s *Storage) GetTransaction(hash common.Hash) (*models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("transaction").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}
