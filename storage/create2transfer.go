package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v4"
)

var (
	create2TransferColumns = []string{
		"transaction_base.*",
		"create2transfer.to_state_id",
		"create2transfer.to_public_key",
	}
	create2TransferWithBatchColumns = []string{
		"transaction_base.*",
		"create2transfer.to_state_id",
		"create2transfer.to_public_key",
		"batch.batch_hash",
		"batch.submission_time",
	}
)

func (s *TransactionStorage) AddCreate2Transfer(t *models.Create2Transfer) (receiveTime *models.Timestamp, err error) {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Postgres: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	receiveTime, err = txStorage.addTransactionBase(&t.TransactionBase, txtype.Create2Transfer)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	_, err = txStorage.database.Postgres.Query(
		txStorage.database.QB.Insert("create2transfer").
			Values(
				t.Hash,
				t.ToStateID,
				t.ToPublicKey,
			),
	).Exec()
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return receiveTime, nil
}

// BatchAddCreate2Transfer contrary to the AddCreate2Transfer method does not set receive_time column on added transfers
func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	tx, txStorage, err := s.BeginTransaction(TxOptions{Postgres: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	txBases := make([]models.TransactionBase, 0, len(txs))
	for i := range txs {
		txBases = append(txBases, txs[i].TransactionBase)
	}
	err = txStorage.BatchAddTransactionBase(txBases)
	if err != nil {
		return err
	}

	query := s.database.QB.Insert("create2transfer")
	for i := range txs {
		query = query.Values(
			txs[i].Hash,
			txs[i].ToStateID,
			txs[i].ToPublicKey,
		)
	}
	res, err := txStorage.database.Postgres.Query(query).Exec()
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRowsAffected
	}
	return tx.Commit()
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select(create2TransferColumns...).
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

func (s *TransactionStorage) GetPendingCreate2Transfers(limit uint32) ([]models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, limit)
	err := s.database.Postgres.Query(
		s.database.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"included_in_commitment": nil, "error_message": nil}).
			OrderBy("transaction_base.nonce ASC", "transaction_base.tx_hash ASC").
			Limit(uint64(limit)),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id int32) ([]models.Create2TransferForCommitment, error) {
	res := make([]models.Create2TransferForCommitment, 0, 32)
	err := s.database.Postgres.Query(
		s.database.QB.Select("transaction_base.tx_hash",
			"transaction_base.from_state_id",
			"transaction_base.amount",
			"transaction_base.fee",
			"transaction_base.nonce",
			"transaction_base.signature",
			"transaction_base.receive_time",
			"create2transfer.to_state_id",
			"create2transfer.to_public_key").
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"included_in_commitment": id}),
	).Into(&res)
	return res, err
}

func (s *TransactionStorage) SetCreate2TransferToStateID(txHash common.Hash, toStateID uint32) error {
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("create2transfer").
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

func (s *Storage) GetCreate2TransfersByPublicKey(publicKey *models.PublicKey) ([]models.Create2TransferWithBatchDetails, error) {
	accounts, err := s.AccountTree.Leaves(publicKey)
	if err != nil {
		return nil, err
	}

	pubKeyIDs := utils.ValueToInterfaceSlice(accounts, "PubKeyID")

	leaves := make([]models.FlatStateLeaf, 0, 1)
	err = s.database.Badger.Find(&leaves, bh.Where("PubKeyID").In(pubKeyIDs...).Index("PubKeyID"))
	if err != nil {
		return nil, err
	}

	stateIDs := make([]uint32, 0, 1)
	for i := range leaves {
		stateIDs = append(stateIDs, leaves[i].StateID)
	}

	res := make([]models.Create2TransferWithBatchDetails, 0, 1)
	err = s.database.Postgres.Query(
		s.database.QB.Select(create2TransferWithBatchColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			LeftJoin("commitment on commitment.commitment_id = transaction_base.included_in_commitment").
			LeftJoin("batch on batch.batch_id = commitment.included_in_batch").
			Where(squirrel.Or{
				squirrel.Eq{"transaction_base.from_state_id": stateIDs},
				squirrel.Eq{"create2transfer.to_state_id": stateIDs},
			}),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetCreate2TransferWithBatchDetails(hash common.Hash) (*models.Create2TransferWithBatchDetails, error) {
	res := make([]models.Create2TransferWithBatchDetails, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select(create2TransferWithBatchColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			LeftJoin("commitment on commitment.commitment_id = transaction_base.included_in_commitment").
			LeftJoin("batch on batch.batch_id = commitment.included_in_batch").
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
