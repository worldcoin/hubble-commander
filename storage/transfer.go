package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var transferColumns = []string{
	"transaction_base.*",
	"transfer.to_state_id",
}

func (s *TransactionStorage) AddTransfer(t *models.Transfer) error {
	t.SetReceiveTime()
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTransactionFromTransfer(t))
}

// BatchAddTransfer contrary to the AddTransfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddTransfer(txs []models.Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		err = txStorage.database.Badger.Insert(txs[i].Hash, models.MakeStoredTransactionFromTransfer(&txs[i]))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *TransactionStorage) GetTransfer(hash common.Hash) (*models.Transfer, error) {
	tx, err := s.getStoredTransaction(hash)
	if err != nil {
		return nil, err
	}
	return tx.ToTransfer(), nil
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
	res := make([]models.StoredTransaction, 0, 32)
	err := s.database.Badger.Find(
		&res,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID"),
	)
	if err != nil {
		return nil, err
	}

	txs := make([]models.TransferForCommitment, 0, len(res))
	for i := range res {
		txs = append(txs, *res[i].ToTransferForCommitment())
	}
	return txs, nil
}

func (s *TransactionStorage) BatchMarkTransfersAsIncluded(txs []models.Transfer, commitmentID *models.CommitmentID) error {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		storedTx := models.MakeStoredTransactionFromTransfer(&txs[i])
		storedTx.CommitmentID = commitmentID
		err = txStorage.updateStoredTransaction(&storedTx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) GetTransferWithBatchDetails(hash common.Hash) (*models.TransferWithBatchDetails, error) {
	transfer, err := s.GetTransfer(hash)
	if err != nil {
		return nil, err
	}
	res := &models.TransferWithBatchDetails{Transfer: *transfer}
	if res.CommitmentID == nil {
		return res, nil
	}

	batch, err := s.GetBatch(res.CommitmentID.BatchID)
	if err != nil {
		return nil, err
	}
	res.BatchHash = batch.Hash
	res.BatchTime = batch.SubmissionTime

	return res, nil
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
