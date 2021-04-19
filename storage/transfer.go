package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var transferColumns = []string{
	"transaction_base.*",
	"transfer.to_state_id",
}

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
	if err != nil {
		return err
	}

	_, err = s.DB.Query(
		s.QB.Insert("transfer").
			Values(
				t.Hash,
				t.ToStateID,
			),
	).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.DB.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transfer")
	}
	return &res[0], nil
}

func (s *Storage) GetUserTransfers(fromStateID models.Uint256) ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.DB.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"from_state_id": fromStateID}),
	).Into(&res)
	return res, err
}

func (s *Storage) GetPendingTransfers() ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, 32)
	err := s.DB.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}), // TODO order by nonce asc, then order by fee desc
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.DB.Query(
		s.QB.Select(transferColumns...).
			From("account").
			JoinClause("NATURAL JOIN state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Join("transaction_base on transaction_base.from_state_id::bit(33) = state_node.state_id").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"account.public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
