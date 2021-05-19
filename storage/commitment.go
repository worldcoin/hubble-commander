package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var commitmentWithTokenIndexCols = []string{
	"commitment.commitment_id",
	"commitment.transactions",
	"commitment.fee_receiver",
	"commitment.combined_signature",
	"commitment.post_state_root",
}

func (s *Storage) AddCommitment(commitment *models.Commitment) (*int32, error) {
	res := make([]int32, 0, 1)
	err := s.Postgres.Query(
		s.QB.Insert("commitment").
			Values(
				squirrel.Expr("DEFAULT"),
				commitment.Type,
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
	err := s.Postgres.Query(
		s.QB.Select("*").
			From("commitment").
			Where(squirrel.Eq{"commitment_id": id}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("commitment")
	}
	return &res[0], nil
}

func (s *Storage) MarkCommitmentAsIncluded(id int32, batchHash, accountRoot *common.Hash) error {
	res, err := s.Postgres.Query(
		s.QB.Update("commitment").
			Where(squirrel.Eq{"commitment_id": id}).
			Set("included_in_batch", *batchHash).
			Set("account_tree_root", *accountRoot),
	).Exec()
	if err != nil {
		return err
	}

	numUpdatedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if numUpdatedRows == 0 {
		return ErrNoRowsAffected
	}
	return nil
}

func (s *Storage) GetPendingCommitments(maxFetched uint64) ([]models.Commitment, error) {
	res := make([]models.Commitment, 0, 32)
	err := s.Postgres.Query(
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

func (s *Storage) GetCommitmentsByBatchHash(hash *common.Hash) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.CommitmentWithTokenID, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select(commitmentWithTokenIndexCols...).
			From("commitment").
			Where(squirrel.Eq{"commitment.included_in_batch": hash}),
	).Into(&commitments)
	if err != nil {
		return nil, err
	}
	if len(commitments) == 0 {
		return nil, NewNotFoundError("commitments")
	}

	for i := range commitments {
		path := models.MerklePath{
			Path:  commitments[i].FeeReceiverStateID,
			Depth: 32,
		}

		var node models.StateNode
		err := s.Badger.Get(path, &node)
		if err != nil {
			return nil, err
		}

		leaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.Badger.Find(&leaves, bh.Where("DataHash").Eq(node.DataHash).Index("DataHash"))
		if err != nil {
			return nil, err
		}
		if len(leaves) != 1 {
			return nil, NewNotFoundError("commitments")
		}

		commitments[i].TokenID = leaves[0].TokenIndex
	}

	return commitments, nil
}

func (s *Storage) GetCommitmentsByBatchID(id models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.CommitmentWithTokenID, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select(commitmentWithTokenIndexCols...).
			From("batch").
			Join("commitment ON commitment.included_in_batch = batch.batch_hash").
			Where(squirrel.Eq{"batch.batch_id": id}),
	).Into(&commitments)
	if err != nil {
		return nil, err
	}
	if len(commitments) == 0 {
		return nil, NewNotFoundError("commitments") // TODO TEST ME
	}

	for i := range commitments {
		path := models.MerklePath{
			Path:  commitments[i].FeeReceiverStateID,
			Depth: 32,
		}

		var node models.StateNode
		err := s.Badger.Get(path, &node)
		if err != nil {
			return nil, err
		}

		leaves := make([]models.FlatStateLeaf, 0, 1)
		err = s.Badger.Find(&leaves, bh.Where("DataHash").Eq(node.DataHash).Index("DataHash"))
		if err != nil {
			return nil, err
		}
		if len(leaves) != 1 {
			return nil, NewNotFoundError("commitments")
		}

		commitments[i].TokenID = leaves[0].TokenIndex
	}

	return commitments, nil
}
