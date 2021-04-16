package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var create2transferColumns = []string{
	"transaction_base.tx_hash",
	"transaction_base.from_state_id",
	"transaction_base.amount",
	"transaction_base.fee",
	"transaction_base.nonce",
	"transaction_base.signature",
	"transaction_base.included_in_commitment",
	"transaction_base.error_message",
	"create2transfer.to_state_id",
	"create2transfer.to_pubkey_id",
}

func (s *Storage) AddCreate2Transfer(t *models.Create2Transfer) error {
	_, err := s.DB.Query(
		s.QB.Insert("transaction_base").
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

	_, err = s.DB.Query(
		s.QB.Insert("create2transfer").
			Values(
				t.Hash,
				t.ToStateID,
				t.ToPubkeyID,
			),
	).Exec()
	if err != nil {
		return err
	}

	return nil
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
