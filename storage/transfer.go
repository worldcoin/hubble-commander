package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var transferColumns = []string{
	"transaction_base.*",
	"transfer.to_state_id",
}

func (s *Storage) AddTransfer(t *models.Transfer) error {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Postgres: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = txStorage.Postgres.Query(
		txStorage.QB.Insert("transaction_base").
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

	_, err = txStorage.Postgres.Query(
		txStorage.QB.Insert("transfer").
			Values(
				t.Hash,
				t.ToStateID,
			),
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"tx_hash": hash}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, NewNotFoundError("transaction")
	}
	return &res[0], nil
}

func (s *Storage) GetUserTransfers(fromStateID models.Uint256) ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"from_state_id": fromStateID}),
	).Into(&res)
	return res, err
}

func (s *Storage) GetPendingTransfers(maxFetched uint64) ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, maxFetched)
	err := s.Postgres.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}).
			OrderBy("transaction_base.nonce ASC").
			Limit(maxFetched),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.Transfer, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(&leaves, bh.Where("PubKeyID").In(pubKeyIDs...).Index("PubKeyID"))
	if err != nil {
		return nil, err
	}

	stateIDs := make([]uint32, 0, 1)
	for i := range leaves {
		stateIDs = append(stateIDs, leaves[i].StateID)
	}

	res := make([]models.Transfer, 0, 1)
	err = s.Postgres.Query(
		s.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Or{
				squirrel.Eq{"transaction_base.from_state_id": stateIDs},
				squirrel.Eq{"transfer.to_state_id": stateIDs},
			}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetTransfersByCommitmentID(id int32) ([]models.TransferForCommitment, error) {
	res := make([]models.TransferForCommitment, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select("transaction_base.tx_hash",
			"transaction_base.from_state_id",
			"transaction_base.amount",
			"transaction_base.fee",
			"transaction_base.nonce",
			"transaction_base.signature",
			"transfer.to_state_id").
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"included_in_commitment": id}),
	).Into(&res)
	return res, err
}
