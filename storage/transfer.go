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

func (s *TransactionStorage) AddTransfer(t *models.Transfer) (receiveTime *models.Timestamp, err error) {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Postgres: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	receiveTime, err = txStorage.addTransactionBase(&t.TransactionBase, txtype.Transfer)
	if err != nil {
		return nil, err
	}

	_, err = txStorage.database.Postgres.Query(
		txStorage.database.QB.Insert("transfer").
			Values(
				t.Hash,
				t.ToStateID,
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

// BatchAddTransfer contrary to the AddTransfer method does not set receive_time column on added transfers
func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
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

	query := s.database.QB.Insert("transfer")
	for i := range txs {
		query = query.Values(
			txs[i].Hash,
			txs[i].ToStateID,
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

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select(transferColumns...).
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

func (s *TransactionStorage) GetPendingTransfers(limit uint32) ([]models.Transfer, error) {
	res := make([]models.Transfer, 0, limit)
	err := s.database.Postgres.Query(
		s.database.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"batch_id": nil, "error_message": nil}).
			OrderBy("transaction_base.nonce ASC", "transaction_base.tx_hash ASC").
			Limit(uint64(limit)),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TransactionStorage) GetTransfersByCommitmentID(id *models.CommitmentID) ([]models.TransferForCommitment, error) {
	res := make([]models.TransferForCommitment, 0, 32)
	err := s.database.Postgres.Query(
		s.database.QB.Select("transaction_base.tx_hash",
			"transaction_base.from_state_id",
			"transaction_base.amount",
			"transaction_base.fee",
			"transaction_base.nonce",
			"transaction_base.signature",
			"transaction_base.receive_time",
			"transfer.to_state_id").
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Eq{"batch_id": id.BatchID, "index_in_batch": id.IndexInBatch}),
	).Into(&res)
	return res, err
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	res := make([]models.Transfer, 0, 1)
	err := s.database.Postgres.Query(
		s.database.QB.Select(transferColumns...).
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
	transfer := &models.TransferWithBatchDetails{Transfer: res[0]}
	if transfer.CommitmentID == nil {
		return transfer, nil
	}

	batch, err := s.GetBatch(transfer.CommitmentID.BatchID)
	if err != nil {
		return nil, err
	}
	transfer.BatchHash = batch.Hash
	transfer.BatchTime = batch.SubmissionTime

	return transfer, nil
}

func (s *Storage) GetTransfersByPublicKey(publicKey *models.PublicKey) ([]models.TransferWithBatchDetails, error) {
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

	transfers := make([]models.Transfer, 0, 1)
	err = s.database.Postgres.Query(
		s.database.QB.Select(transferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN transfer").
			Where(squirrel.Or{
				squirrel.Eq{"transaction_base.from_state_id": stateIDs},
				squirrel.Eq{"transfer.to_state_id": stateIDs},
			}),
	).Into(&transfers)
	if err != nil {
		return nil, err
	}

	return s.transfersToTransfersWithBatchDetails(transfers)
}

func (s *Storage) transfersToTransfersWithBatchDetails(transfers []models.Transfer) (result []models.TransferWithBatchDetails, err error) {
	result = make([]models.TransferWithBatchDetails, 0, len(transfers))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range transfers {
		if transfers[i].CommitmentID == nil {
			result = append(result, models.TransferWithBatchDetails{Transfer: transfers[i]})
			continue
		}
		batch, ok := batchIDs[transfers[i].CommitmentID.BatchID]
		if !ok {
			batch, err = s.GetBatch(transfers[i].CommitmentID.BatchID)
			if err != nil {
				return nil, err
			}
			batchIDs[transfers[i].CommitmentID.BatchID] = batch
		}

		result = append(result, models.TransferWithBatchDetails{
			Transfer:  transfers[i],
			BatchHash: batch.Hash,
			BatchTime: batch.SubmissionTime,
		})
	}
	return result, nil
}
