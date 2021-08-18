package storage

import (
	"github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/db/badger"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	bh "github.com/timshannon/badgerhold/v3"
)

type CommitmentStorage struct {
	database *Database
}

func NewCommitmentStorage(database *Database) *CommitmentStorage {
	return &CommitmentStorage{
		database: database,
	}
}

func (s *CommitmentStorage) copyWithNewDatabase(database *Database) *CommitmentStorage {
	newCommitmentStorage := *s
	newCommitmentStorage.database = database

	return &newCommitmentStorage
}

func (s *CommitmentStorage) AddCommitment(commitment *models.Commitment) error {
	err := s.database.Badger.Insert(models.CommitmentKey{
		BatchID:      commitment.BatchID,
		IndexInBatch: commitment.IndexInBatch,
	}, *commitment)
	return err
}

func (s *CommitmentStorage) GetCommitment(batchID models.Uint256, commitmentIndex uint32) (*models.Commitment, error) {
	commitment := models.Commitment{
		BatchID:      batchID,
		IndexInBatch: commitmentIndex,
	}
	err := s.database.Badger.Get(models.CommitmentKey{
		BatchID:      batchID,
		IndexInBatch: commitmentIndex,
	}, &commitment)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("commitment")
	}
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func (s *CommitmentStorage) GetLatestCommitment() (*models.Commitment, error) {
	var commitment *models.Commitment
	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		opts := bdg.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Reverse = true
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(models.CommitmentPrefix)+1)
		seekPrefix = append(seekPrefix, models.CommitmentPrefix...)
		seekPrefix = append(seekPrefix, 0xFF) // Required to loop backwards

		it.Seek(seekPrefix)
		if it.ValidForPrefix(models.CommitmentPrefix) {
			var err error
			commitment, err = decodeCommitment(it.Item())
			return err
		}
		return NewNotFoundError("commitment")
	})
	if err != nil {
		return nil, err
	}

	return commitment, nil
}

func (s *CommitmentStorage) MarkCommitmentAsIncluded(commitmentID int32, batchID models.Uint256) error {
	//TODO-dis: remove
	res, err := s.database.Postgres.Query(
		s.database.QB.Update("commitment").
			Where(squirrel.Eq{"commitment_id": commitmentID}).
			Set("included_in_batch", batchID),
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

func (s *CommitmentStorage) DeleteCommitmentsByBatchIDs(batchIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	txn, err := txDatabase.Badger.Txn()
	if err != nil {
		return err
	}

	keys := make([]models.CommitmentKey, 0, len(batchIDs))
	opts := bdg.DefaultIteratorOptions
	opts.PrefetchValues = true
	for i := range batchIDs {
		commitmentKeys, err := getCommitmentKeysByBatchID(txn, opts, batchIDs[i])
		if err != nil {
			return err
		}
		keys = append(keys, commitmentKeys...)
	}

	var commitment models.Commitment
	for i := range keys {
		err = txDatabase.Badger.Delete(keys[i], commitment)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func getCommitmentKeysByBatchID(txn *bdg.Txn, opts bdg.IteratorOptions, batchID models.Uint256) ([]models.CommitmentKey, error) {
	keys := make([]models.CommitmentKey, 0, 32)
	it := txn.NewIterator(opts)
	defer it.Close()

	seekPrefix := make([]byte, 0, len(models.CommitmentPrefix)+32)
	seekPrefix = append(seekPrefix, models.CommitmentPrefix...)
	seekPrefix = append(seekPrefix, utils.PadLeft(batchID.Bytes(), 32)...)

	for it.Seek(seekPrefix); it.ValidForPrefix(seekPrefix); it.Next() {
		var key models.CommitmentKey
		err := key.SetBytes(it.Item().Key()[len(models.CommitmentPrefix):])
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.Commitment, 0, 32)
	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		opts := bdg.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(models.CommitmentPrefix)+33)
		seekPrefix = append(seekPrefix, models.CommitmentPrefix...)
		seekPrefix = append(seekPrefix, utils.PadLeft(batchID.Bytes(), 32)...)

		for it.Seek(seekPrefix); it.ValidForPrefix(seekPrefix); it.Next() {
			commitment, err := decodeCommitment(it.Item())
			if err != nil {
				return err
			}
			commitments = append(commitments, *commitment)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(commitments) == 0 {
		return nil, NewNotFoundError("commitments")
	}

	commitmentsWithToken := make([]models.CommitmentWithTokenID, 0, len(commitments))
	for i := range commitments {
		stateLeaf, err := s.StateTree.Leaf(commitments[i].FeeReceiver)
		if err != nil {
			return nil, err
		}
		commitmentsWithToken = append(commitmentsWithToken, models.CommitmentWithTokenID{
			ID:                 0,
			Transactions:       commitments[i].Transactions,
			TokenID:            stateLeaf.TokenID,
			FeeReceiverStateID: commitments[i].FeeReceiver,
			CombinedSignature:  commitments[i].CombinedSignature,
			PostStateRoot:      commitments[i].PostStateRoot,
		})
	}

	return commitmentsWithToken, nil
}

func decodeCommitment(item *bdg.Item) (*models.Commitment, error) {
	var commitment models.Commitment
	err := item.Value(func(v []byte) error {
		return badger.Decode(v, &commitment)
	})
	if err != nil {
		return nil, err
	}

	var key models.CommitmentKey
	err = badger.DecodeKey(item.Key(), &key, models.CommitmentPrefix)
	if err != nil {
		return nil, err
	}
	commitment.BatchID = key.BatchID
	commitment.IndexInBatch = key.IndexInBatch
	return &commitment, nil
}
