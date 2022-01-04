package storage

import (
	"reflect"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

const (
	AccountTreeDepth = merkletree.MaxDepth

	AccountBatchSize     = 1 << 4
	AccountBatchOffset   = 1 << 31
	leftSubtreeMaxValue  = AccountBatchOffset - 2
	rightSubtreeMaxValue = AccountBatchOffset*2 - 18
)

var (
	ErrInvalidAccountsLength = errors.New("invalid accounts length")
	accountLeafPrefix        = []byte("bh_" + reflect.TypeOf(models.AccountLeaf{}).Name())
)

type AccountTree struct {
	database   *Database
	merkleTree *StoredMerkleTree
}

func NewAccountTree(database *Database) *AccountTree {
	return &AccountTree{
		database:   database,
		merkleTree: NewStoredMerkleTree("account", database, AccountTreeDepth),
	}
}

func (s *AccountTree) copyWithNewDatabase(database *Database) *AccountTree {
	return NewAccountTree(database)
}

func (s *AccountTree) Root() (*common.Hash, error) {
	return s.merkleTree.Root()
}

func (s *AccountTree) Leaf(pubKeyID uint32) (*models.AccountLeaf, error) {
	var leaf models.AccountLeaf
	err := s.database.Badger.Get(pubKeyID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("account leaf")
	}
	if err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (s *AccountTree) Leaves(publicKey *models.PublicKey) ([]models.AccountLeaf, error) {
	accounts := make([]models.AccountLeaf, 0, 1)
	err := s.database.Badger.Find(
		&accounts,
		bh.Where("PublicKey").Eq(*publicKey).Index("PublicKey"),
	)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, NewNotFoundError("account leaves")
	}
	return accounts, nil
}

func (s *AccountTree) SetSingle(leaf *models.AccountLeaf) error {
	if leaf.PubKeyID > leftSubtreeMaxValue {
		return NewInvalidPubKeyIDError(leaf.PubKeyID)
	}

	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = NewAccountTree(txDatabase).unsafeSet(leaf)
	if err == bh.ErrKeyExists {
		return NewAccountAlreadyExistsError(leaf)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AccountTree) SetBatch(leaves []models.AccountLeaf) error {
	if len(leaves) != AccountBatchSize {
		return ErrInvalidAccountsLength
	}

	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	accountTree := NewAccountTree(txDatabase)

	for i := range leaves {
		if leaves[i].PubKeyID < AccountBatchOffset || leaves[i].PubKeyID > rightSubtreeMaxValue {
			return NewInvalidPubKeyIDError(leaves[i].PubKeyID)
		}
		_, err = accountTree.unsafeSet(&leaves[i])
		if err == bh.ErrKeyExists {
			return NewAccountBatchAlreadyExistsError(leaves)
		}
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *AccountTree) GetWitness(pubKeyID uint32) (models.Witness, error) {
	return s.merkleTree.GetWitness(models.MakeMerklePathFromLeafID(pubKeyID))
}

func (s *AccountTree) unsafeSet(leaf *models.AccountLeaf) (models.Witness, error) {
	err := s.database.Badger.Insert(leaf.PubKeyID, *leaf)
	if err != nil {
		return nil, err
	}

	dataHash := crypto.Keccak256Hash(leaf.PublicKey.Bytes())
	path := models.MakeMerklePathFromLeafID(leaf.PubKeyID)
	_, witness, err := s.merkleTree.SetNode(&path, dataHash)
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *AccountTree) NextBatchAccountPubKeyID() (*uint32, error) {
	nextPubKeyID := uint32(AccountBatchOffset)
	err := s.database.Badger.View(func(txn *bdg.Txn) error {
		opts := bdg.IteratorOptions{
			PrefetchValues: false,
			Reverse:        true,
		}
		it := txn.NewIterator(opts)
		defer it.Close()

		seekPrefix := make([]byte, 0, len(accountLeafPrefix)+1)
		seekPrefix = append(seekPrefix, accountLeafPrefix...)
		seekPrefix = append(seekPrefix, 0xFF)
		it.Seek(seekPrefix)
		if !it.ValidForPrefix(accountLeafPrefix) {
			return nil
		}

		var account models.AccountLeaf
		err := it.Item().Value(account.SetBytes)
		if err != nil {
			return err
		}

		if account.PubKeyID < AccountBatchOffset {
			return nil
		}
		nextPubKeyID = account.PubKeyID + 1
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &nextPubKeyID, nil
}
