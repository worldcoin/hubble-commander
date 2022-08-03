package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (s *TransactionStorage) GetMassMigration(hash common.Hash) (*models.MassMigration, error) {
	tx, err := s.getTransactionByHash(hash)
	if err != nil {
		return nil, err
	}
	if tx.Type() != txtype.MassMigration {
		return nil, errors.WithStack(NewNotFoundError("transaction"))
	}
	transfer := tx.ToMassMigration()
	return transfer, nil
}
