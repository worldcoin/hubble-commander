package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var create2transferColumns = []string{
	"transaction_base.*",
	"create2transfer.to_state_id",
	"create2transfer.to_pub_key_id",
}

func (s *Storage) AddCreate2Transfer(t *models.Create2Transfer) (err error) {
	tx, txStorage, err := s.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = txStorage.DB.Query(
		txStorage.QB.Insert("transaction_base").
			Values(
				t.Hash,
				txtype.Create2Transfer,
				t.FromStateID,
				t.Amount,
				t.Fee,
				t.Nonce,
				t.Signature,
				t.IncludedInCommitment,
				t.ErrorMessage,
			),
	).Exec()
	if err != nil {
		return err
	}

	_, err = txStorage.DB.Query(
		txStorage.QB.Insert("create2transfer").
			Values(
				t.Hash,
				t.ToStateID,
				t.ToPubKeyID,
			),
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, 1)
	err := s.DB.Query(
		s.QB.Select(create2transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("create2transfer")
	}
	return &res[0], nil
}
