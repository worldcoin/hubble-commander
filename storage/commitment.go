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
	return s.database.Badger.Insert(commitment.ID, *commitment)
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
	var err error
	err = s.database.Badger.ReverseIterator(models.CommitmentPrefix, func(item *bdg.Item) (bool, error) {
		commitment, err = decodeCommitment(item)
		return true, err
	})
	if err == badger.ErrIteratorFinished {
		return nil, NewNotFoundError("commitment")
	}
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

	var commitmentIDs []models.CommitmentID
	opts := bdg.DefaultIteratorOptions
	opts.PrefetchValues = false
	ids := make([]models.CommitmentID, 0, len(batchIDs))
	for i := range batchIDs {
		commitmentIDs, err = getCommitmentIDsByBatchID(txDatabase, opts, batchIDs[i])
		if err != nil {
			return err
		}
		ids = append(ids, commitmentIDs...)
	}

	if len(ids) == 0 {
		return NewNotFoundError("commitments")
	}

	var commitment models.Commitment
	for i := range ids {
		err = txDatabase.Badger.Delete(ids[i], commitment)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func getCommitmentIDsByBatchID(txn *Database, opts bdg.IteratorOptions, batchID models.Uint256) ([]models.CommitmentID, error) {
	ids := make([]models.CommitmentID, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := txn.Badger.Iterator(prefix, opts, func(item *bdg.Item) (bool, error) {
		var id models.CommitmentID
		err := badger.DecodeKey(item.Key(), &id, models.CommitmentPrefix)
		if err != nil {
			return false, err
		}
		ids = append(ids, id)
		return false, nil
	})
	if err != nil && err != badger.ErrIteratorFinished {
		return nil, err
	}
	return ids, nil
}

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.Commitment, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := s.database.Badger.Iterator(prefix, bdg.DefaultIteratorOptions, func(item *bdg.Item) (bool, error) {
		commitment, err := decodeCommitment(item)
		if err != nil {
			return false, err
		}
		commitments = append(commitments, *commitment)
		return false, nil
	})
	if err != nil && err != badger.ErrIteratorFinished {
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
			ID:                 commitments[i].ID,
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

func getCommitmentPrefixWithBatchID(batchID *models.Uint256) []byte {
	commitmentPrefixLen := len(models.CommitmentPrefix)
	prefix := make([]byte, commitmentPrefixLen+32)
	copy(prefix[:commitmentPrefixLen], models.CommitmentPrefix)
	copy(prefix[commitmentPrefixLen:], utils.PadLeft(batchID.Bytes(), 32))
	return prefix
}
