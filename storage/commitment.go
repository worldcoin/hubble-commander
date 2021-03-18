package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (s *Storage) AddCommitment(commitment *models.Commitment) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Insert("commitment").
			Values(
				commitment.LeafHash,
				commitment.PostStateRoot,
				commitment.BodyHash,
				commitment.AccountTreeRoot,
				commitment.CombinedSignature,
				commitment.FeeReceiver,
				commitment.Transactions,
				commitment.IncludedInBatch,
			),
	)
	return err
}

func (s *Storage) GetCommitment(leafHash common.Hash) (*models.Commitment, error) {
	res := make([]models.Commitment, 0, 1)
	err := s.DB.Query(
		squirrel.Select("*").
			From("commitment").
			Where(squirrel.Eq{"leaf_hash": leafHash}),
	).Into(&res)
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return &res[0], nil
}

func (s *Storage) MarkCommitmentAsIncluded(leafHash, batchHash common.Hash) error {
	_, err := s.DB.ExecBuilder(
		s.QB.Update("commitment").
			Where(squirrel.Eq{"leaf_hash": leafHash}).
			Set("included_in_batch", batchHash),
	)
	return err
}
