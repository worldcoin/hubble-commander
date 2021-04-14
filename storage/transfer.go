package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddTransfer(t *models.Transfer) error {
	_, err := s.DB.Query(
		s.QB.Insert("transaction_base").
			Values(
				t.Hash,
				txtype.Transfer,
				t.FromStateID,
				t.Amount,
				t.Fee,
				t.Nonce,
				t.Signature,
				t.IncludedInCommitment,
				t.ErrorMessage,
			),
	).Exec()

	_, err = s.DB.Query(
		s.QB.Insert("transfer").
			Values(
				t.Hash,
				t.ToStateID,
			),
	).Exec()

	return err
}

func (s *Storage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.DB.Query(
		s.QB.Select(
			"transaction_base.tx_hash",
			"transaction_base.from_state_id",
			"transaction_base.amount",
			"transaction_base.fee",
			"transaction_base.nonce",
			"transaction_base.signature",
			"transaction_base.included_in_commitment",
			"transaction_base.error_message",
			"transfer.to_state_id",
		).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}
