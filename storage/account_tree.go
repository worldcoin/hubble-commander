package storage

import (
	"github.com/Worldcoin/hubble-commander/db"
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

var ErrInvalidAccountsLength = errors.New("invalid accounts length")

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
		return nil, errors.WithStack(NewNotFoundError("account leaf"))
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
		bh.Where("PublicKey").Eq(publicKey).Index("PublicKey"),
	)
	if err != nil {
		return nil, err
	}
	if len(accounts) == 0 {
		return nil, errors.WithStack(NewNotFoundError("account leaves"))
	}
	return accounts, nil
}

func (s *AccountTree) SetSingle(leaf *models.AccountLeaf) error {
	if leaf.PubKeyID > leftSubtreeMaxValue {
		return errors.WithStack(NewInvalidPubKeyIDError(leaf.PubKeyID))
	}

	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = NewAccountTree(txDatabase).unsafeSet(leaf)
	if err == bh.ErrKeyExists {
		return errors.WithStack(NewAccountAlreadyExistsError(leaf))
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

	return s.SetInBatch(leaves...)
}

func (s *AccountTree) SetInBatch(leaves ...models.AccountLeaf) error {
	tx, txDatabase, err := s.database.BeginTransaction(TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	accountTree := NewAccountTree(txDatabase)

	for i := range leaves {
		if isValidBatchAccount(&leaves[i]) {
			return errors.WithStack(NewInvalidPubKeyIDError(leaves[i].PubKeyID))
		}
		_, err = accountTree.unsafeSet(&leaves[i])
		if err == bh.ErrKeyExists {
			return errors.WithStack(NewAccountBatchAlreadyExistsError(leaves))
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
	err := s.database.Badger.Iterator(models.AccountLeafPrefix, db.ReverseKeyIteratorOpts, func(item *bdg.Item) (finish bool, err error) {
		var account models.AccountLeaf
		err = item.Value(account.SetBytes)
		if err != nil {
			return false, err
		}

		if account.PubKeyID < AccountBatchOffset {
			return true, nil
		}
		nextPubKeyID = account.PubKeyID + 1
		return true, nil
	})
	if err != nil && err != db.ErrIteratorFinished {
		return nil, err
	}
	return &nextPubKeyID, nil
}

func isValidBatchAccount(leaf *models.AccountLeaf) bool {
	return leaf.PubKeyID < AccountBatchOffset || leaf.PubKeyID > rightSubtreeMaxValue
}
