package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

var create2TransferColumns = []string{
	"transaction_base.*",
	"create2transfer.to_state_id",
	"create2transfer.to_public_key",
}

func (s *TransactionStorage) AddCreate2Transfer(t *models.Create2Transfer) error {
	t.SetReceiveTime()
	return s.database.Badger.Insert(t.Hash, models.MakeStoredTransactionFromCreate2Transfer(t))
}

// BatchAddCreate2Transfer contrary to the AddCreate2Transfer method does not set ReceiveTime field on added transfers
func (s *TransactionStorage) BatchAddCreate2Transfer(txs []models.Create2Transfer) error {
	if len(txs) < 1 {
		return ErrNoRowsAffected
	}

	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		err = txStorage.database.Badger.Insert(txs[i].Hash, models.MakeStoredTransactionFromCreate2Transfer(&txs[i]))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *TransactionStorage) GetCreate2Transfer(hash common.Hash) (*models.Create2Transfer, error) {
	tx, err := s.getStoredTransaction(hash)
	if err != nil {
		return nil, err
	}
	return tx.ToCreate2Transfer(), nil
}

func (s *TransactionStorage) GetPendingCreate2Transfers(limit uint32) ([]models.Create2Transfer, error) {
	res := make([]models.Create2Transfer, 0, limit)
	err := s.database.Postgres.Query(
		s.database.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Eq{"batch_id": nil, "error_message": nil}).
			OrderBy("transaction_base.nonce ASC", "transaction_base.tx_hash ASC").
			Limit(uint64(limit)),
	).Into(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *TransactionStorage) GetCreate2TransfersByCommitmentID(id *models.CommitmentID) ([]models.Create2TransferForCommitment, error) {
	res := make([]models.StoredTransaction, 0, 32)
	err := s.database.Badger.Find(
		&res,
		bh.Where("CommitmentID").Eq(*id).Index("CommitmentID"),
	)
	if err != nil {
		return nil, err
	}

	txs := make([]models.Create2TransferForCommitment, 0, len(res))
	for i := range res {
		txs = append(txs, *res[i].ToCreate2TransferForCommitment())
	}
	return txs, nil
}

// SetCreate2TransferToStateID TODO-tx: remove
func (s *TransactionStorage) SetCreate2TransferToStateID(txHash common.Hash, toStateID uint32) error {
	transfer, err := s.GetCreate2Transfer(txHash)
	if err != nil {
		return err
	}
	transfer.ToStateID = ref.Uint32(toStateID)
	err = s.database.Badger.Update(txHash, models.MakeStoredTransactionFromCreate2Transfer(transfer))
	if err == bh.ErrNotFound {
		return NewNotFoundError("transaction")
	}
	return err
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

	transfers := make([]models.Create2Transfer, 0, 1)
	err = s.database.Postgres.Query(
		s.database.QB.Select(create2TransferColumns...).
			From("transaction_base").
			JoinClause("NATURAL JOIN create2transfer").
			Where(squirrel.Or{
				squirrel.Eq{"transaction_base.from_state_id": stateIDs},
				squirrel.Eq{"create2transfer.to_state_id": stateIDs},
			}),
	).Into(&transfers)
	if err != nil {
		return nil, err
	}

	return s.create2TransferToTransfersWithBatchDetails(transfers)
}

func (s *TransactionStorage) BatchMarkCreate2TransfersAsIncluded(txs []models.Create2Transfer, commitmentID *models.CommitmentID) error {
	tx, txStorage, err := s.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	for i := range txs {
		storedTx := models.MakeStoredTransactionFromCreate2Transfer(&txs[i])
		storedTx.CommitmentID = commitmentID
		err = txStorage.updateStoredTransaction(&storedTx)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) GetCreate2TransferWithBatchDetails(hash common.Hash) (*models.Create2TransferWithBatchDetails, error) {
	res, err := s.GetCreate2Transfer(hash)
	if err != nil {
		return nil, err
	}
	transfer := &models.Create2TransferWithBatchDetails{Create2Transfer: *res}
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

func (s *Storage) create2TransferToTransfersWithBatchDetails(transfers []models.Create2Transfer) (
	result []models.Create2TransferWithBatchDetails,
	err error,
) {
	result = make([]models.Create2TransferWithBatchDetails, 0, len(transfers))
	batchIDs := make(map[models.Uint256]*models.Batch)
	for i := range transfers {
		if transfers[i].CommitmentID == nil {
			result = append(result, models.Create2TransferWithBatchDetails{Create2Transfer: transfers[i]})
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

		result = append(result, models.Create2TransferWithBatchDetails{
			Create2Transfer: transfers[i],
			BatchHash:       batch.Hash,
			BatchTime:       batch.SubmissionTime,
		})
	}
	return result, nil
}
