package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddCommitment(commitment *models.Commitment) (*int32, error) {
	res := make([]int32, 0, 1)
	err := s.DB.InsertQuery(
		s.QB.Insert("commitment").
			Values(
				squirrel.Expr("DEFAULT"),
				commitment.Transactions,
				commitment.FeeReceiver,
				commitment.CombinedSignature,
				commitment.PostStateRoot,
				commitment.AccountTreeRoot,
				commitment.IncludedInBatch,
			).
			Suffix("RETURNING commitment_id"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return ref.Int32(res[0]), nil
}

func (s *Storage) GetCommitment(id int32) (*models.Commitment, error) {
	res := make([]models.Commitment, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("commitment").
			Where(squirrel.Eq{"commitment_id": id}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}

func (s *Storage) MarkCommitmentAsIncluded(id int32, batchHash common.Hash) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Update("commitment").
			Where(squirrel.Eq{"commitment_id": id}).
			Set("included_in_batch", batchHash),
	)
	return err
}

func (s *Storage) GetPendingCommitments(maxFetched uint64) ([]models.Commitment, error) {
	res := make([]models.Commitment, 0, 32)
	err := s.DB.Query(
		s.QB.Select("*").
			From("commitment").
			Where(squirrel.Eq{"included_in_batch": nil}).
			Limit(maxFetched),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
