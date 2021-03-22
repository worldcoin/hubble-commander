package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddTransaction(tx *models.Transaction) error {
	_, err := s.DB.ExecBuilder(
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

func (s *Storage) GetPendingTransactions() ([]models.Transaction, error) {
	res := make([]models.Transaction, 0, 32)
	err := s.DB.Query(
		squirrel.Select("*").
			From("transaction").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}), // TODO order by nonce asc, then order by fee desc
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) MarkTransactionAsIncluded(txHash common.Hash, commitmentID int32) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Update("transaction").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("included_in_commitment", commitmentID),
	)
	return err
}

func (s *Storage) SetTransactionError(txHash common.Hash, errorMessage string) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Update("transaction").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("error_message", errorMessage),
	)
	return err
}
