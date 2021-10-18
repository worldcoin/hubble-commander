package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
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

func (s *CommitmentStorage) AddCommitment(commitment *models.TxCommitment) error {
	return s.database.Badger.Insert(commitment.ID, models.MakeStoredCommitmentFromTxCommitment(commitment))
}

func (s *CommitmentStorage) GetCommitment(id *models.CommitmentID) (*models.TxCommitment, error) {
	commitment, err := s.GetStoredCommitment(id)
	if err != nil {
		return nil, err
	}
	return commitment.ToTxCommitment(), nil
}

func (s *CommitmentStorage) DeleteCommitmentsByBatchIDs(batchIDs ...models.Uint256) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	var commitmentIDs []models.CommitmentID
	ids := make([]models.CommitmentID, 0, len(batchIDs))
	for i := range batchIDs {
		commitmentIDs, err = getCommitmentIDsByBatchID(txDatabase, db.ReverseKeyIteratorOpts, batchIDs[i])
		if err != nil {
			return err
		}
		ids = append(ids, commitmentIDs...)
	}

	if len(ids) == 0 {
		return errors.WithStack(NewNotFoundError("commitments"))
	}

	var commitment models.TxCommitment
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
		err := db.DecodeKey(item.Key(), &id, models.TxCommitmentPrefix)
		if err != nil {
			return false, err
		}
		ids = append(ids, id)
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	return ids, nil
}

func (s *Storage) GetCommitmentsByBatchID(batchID models.Uint256) ([]models.CommitmentWithTokenID, error) {
	commitments := make([]models.TxCommitment, 0, 32)
	prefix := getCommitmentPrefixWithBatchID(&batchID)
	err := s.database.Badger.Iterator(prefix, bdg.DefaultIteratorOptions, func(item *bdg.Item) (bool, error) {
		commitment, err := decodeCommitment(item)
		if err != nil {
			return false, err
		}
		commitments = append(commitments, *commitment)
		return false, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
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

func decodeCommitment(item *bdg.Item) (*models.TxCommitment, error) {
	var commitment models.TxCommitment
	err := item.Value(func(v []byte) error {
		return db.Decode(v, &commitment)
	})
	if err != nil {
		return nil, err
	}

	err = db.DecodeKey(item.Key(), &commitment.ID, models.TxCommitmentPrefix)
	if err != nil {
		return nil, err
	}
	return &commitment, nil
}

func getCommitmentPrefixWithBatchID(batchID *models.Uint256) []byte {
	commitmentPrefixLen := len(models.TxCommitmentPrefix)
	prefix := make([]byte, commitmentPrefixLen+32)
	copy(prefix[:commitmentPrefixLen], models.TxCommitmentPrefix)
	copy(prefix[commitmentPrefixLen:], utils.PadLeft(batchID.Bytes(), 32))
	return prefix
}
