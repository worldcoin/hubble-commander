package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := s.DB.Query(
		s.QB.Insert("transaction").
			Values(
				tx.Hash,
				tx.FromIndex,
				tx.ToIndex,
				tx.Amount,
				tx.Fee,
				tx.Nonce,
				tx.Signature,
				tx.IncludedInCommitment,
				tx.ErrorMessage,
			),
	).Exec()

	return err
}

func (s *Storage) GetTransaction(hash common.Hash) (*models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("transaction").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}

func (s *Storage) GetUserTransactions(fromIndex models.Uint256) ([]models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := s.DB.Query(
		s.QB.Select("*").
			From("transaction").
			Where(squirrel.Eq{"from_index": fromIndex}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetPendingTransactions() ([]models.Transaction, error) {
	res := make([]models.Transaction, 0, 32)
	err := s.DB.Query(
		s.QB.Select("*").
			From("transaction").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}), // TODO order by nonce asc, then order by fee desc
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) MarkTransactionAsIncluded(txHash common.Hash, commitmentID int32) error {
	res, err := s.DB.Query(
		s.QB.Update("transaction_base").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("included_in_commitment", commitmentID),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return fmt.Errorf("no rows were affected by the update")
	}
	return nil
}

func (s *Storage) SetTransactionError(txHash common.Hash, errorMessage string) error {
	res, err := s.DB.Query(
		s.QB.Update("transaction").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("error_message", errorMessage),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return fmt.Errorf("no rows were affected by the update")
	}
	return nil
}

func (s *Storage) GetTransactionsByPublicKey(publicKey *models.PublicKey) ([]models.Transaction, error) {
	res := make([]models.Transaction, 0, 1)
	err := s.DB.Query(
		s.QB.Select("transaction.*").
			From("account").
			JoinClause("NATURAL JOIN state_leaf").
			JoinClause("NATURAL JOIN state_node").
			Join("transaction on transaction.from_index::bit(33) = state_node.merkle_path").
			Where(squirrel.Eq{"account.public_key": publicKey}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no transactions found")
	}
	return res, nil
}
