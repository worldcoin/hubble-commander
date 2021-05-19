package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var create2TransferColumns = []string{
	"transaction_base.*",
	"create2transfer.to_state_id",
	"create2transfer.to_public_key",
}

func (s *Storage) AddCreate2Transfer(t *models.Create2Transfer) (err error) {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Postgres: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = txStorage.Postgres.Query(
		txStorage.QB.Insert("transaction_base").
			Values(
				t.Hash,
				txtype.Create2Transfer,
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
		txStorage.QB.Insert("create2transfer").
			Values(
				t.Hash,
				t.ToStateID,
				t.ToPublicKey,
			),
	).Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, 1)
	err := s.Postgres.Query(
		s.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
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

func (s *Storage) GetPendingCreate2Transfers() ([]models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}).
			OrderBy("transaction_base.nonce ASC"),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetCreate2TransfersByPublicKey(publicKey *models.PublicKey) ([]models.Create2Transfer, error) {
	accounts, err := s.GetAccounts(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	// TODO possibly performance killer
	// First get all state nodes and then query leaves by pubkey id and datahash
	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.Badger.Find(&leaves, bh.Where("PubKeyID").In(pubKeyIDs...).Index("PubKeyID"))
	if err != nil {
		return nil, err
	}

	dataHashes := utils.ValueToInterfaceSlice(leaves, "DataHash")

	nodes := make([]models.StateNode, 0, 1)
	err = s.Badger.Find(&nodes, bh.Where("DataHash").In(dataHashes...).Index("DataHash"))
	if err != nil {
		return nil, err
	}

	stateIDs := make([]uint32, 0, 1)
	for i := range nodes {
		stateIDs = append(stateIDs, nodes[i].MerklePath.Path)
	}

	res := make([]models.Create2Transfer, 0, 1)
	err = s.Postgres.Query(
		s.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"transaction_base.from_state_id": stateIDs}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetCreate2TransfersByCommitmentID(id int32) ([]models.Create2TransferForCommitment, error) {
	res := make([]models.Create2TransferForCommitment, 0, 32)
	err := s.Postgres.Query(
		s.QB.Select("transaction_base.tx_hash",
			"transaction_base.from_state_id",
			"transaction_base.amount",
			"transaction_base.fee",
			"transaction_base.nonce",
			"transaction_base.signature",
			"create2transfer.to_state_id",
			"create2transfer.to_public_key").
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"included_in_commitment": id}),
	).Into(&res)
	return res, err
}

func (s *Storage) SetCreate2TransferToStateID(txHash common.Hash, toStateID uint32) error {
	res, err := s.Postgres.Query(
		s.QB.Update("create2transfer").
			Where(squirrel.Eq{"tx_hash": txHash}).
			Set("to_state_id", toStateID),
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
