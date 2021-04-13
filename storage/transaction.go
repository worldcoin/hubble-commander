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
	_, err := s.DB.Query(
		s.QB.Update("transaction").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("included_in_commitment", commitmentID),
	).Exec()
	return err
}

func (s *Storage) SetTransactionError(txHash common.Hash, errorMessage string) error {
	_, err := s.DB.Query(
		s.QB.Update("transaction").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("error_message", errorMessage),
	).Exec()
	return err
}

func (s *Storage) GetTransactions(publicKey *models.PublicKey) ([]models.Transaction, error) {
	//query := `
	//SELECT transaction.*
	//FROM account
	//inner JOIN state_leaf on state_leaf.account_index = account.account_index
	//inner join state_node on state_node.data_hash = state_leaf.data_hash
	//inner join transaction on transaction.from_index::bit(33) = state_node.merkle_path
	//where account.public_key=$1`

	res := make([]models.Transaction, 0, 1)
	err := s.DB.Query(
		s.QB.Select("transaction.*").
			From("account").
			InnerJoin("state_leaf on state_leaf.account_index = account.account_index").
			InnerJoin("state_node on state_node.data_hash = state_leaf.data_hash").
			InnerJoin("transaction on transaction.from_index::bit(33) = state_node.merkle_path").
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
