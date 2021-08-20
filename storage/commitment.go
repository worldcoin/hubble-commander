package storage

import (
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
	err := s.database.Badger.Insert(commitment.ID, *commitment)
	return err
}

func (s *CommitmentStorage) GetCommitment(key *models.CommitmentID) (*models.Commitment, error) {
	commitment := models.Commitment{
		ID: *key,
	}
	err := s.database.Badger.Get(*key, &commitment)
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

	var commitmentIDs []models.CommitmentID
	keys := make([]models.CommitmentID, 0, len(batchIDs))
	for i := range batchIDs {
		commitmentIDs, err = getCommitmentIDsByBatchID(txn, bdg.DefaultIteratorOptions, batchIDs[i])
		if err != nil {
			return err
		}
		keys = append(keys, commitmentIDs...)
	}

	if len(keys) == 0 {
		return ErrNoRowsAffected
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

func getCommitmentIDsByBatchID(txn *bdg.Txn, opts bdg.IteratorOptions, batchID models.Uint256) ([]models.CommitmentID, error) {
	keys := make([]models.CommitmentID, 0, 32)
	it := txn.NewIterator(opts)
	defer it.Close()

	seekPrefix := make([]byte, 0, len(models.CommitmentPrefix)+32)
	seekPrefix = append(seekPrefix, models.CommitmentPrefix...)
	seekPrefix = append(seekPrefix, utils.PadLeft(batchID.Bytes(), 32)...)

	for it.Seek(seekPrefix); it.ValidForPrefix(seekPrefix); it.Next() {
		var key models.CommitmentID
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

	err = badger.DecodeKey(item.Key(), &commitment.ID, models.CommitmentPrefix)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}
