package storage

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) GetLatestTransactionNonce(accountStateID uint32) (*models.Uint256, error) {
	res := make([]models.Uint256, 0, 1)
	err := s.DB.Query(
		s.QB.Select("transaction_base.nonce").
			From("transaction_base").
			Where(squirrel.Eq{"from_state_id": accountStateID}).
			OrderBy("nonce DESC").
			Limit(1),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return &res[0], nil
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
		s.QB.Update("transaction_base").
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
