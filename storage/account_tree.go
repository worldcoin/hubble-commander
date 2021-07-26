package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const (
	AccountTreeDepth = merkletree.MaxDepth

	accountBatchOffset   = 1 << 31
	leftSubtreeMaxValue  = accountBatchOffset - 1
	rightSubtreeMaxValue = accountBatchOffset*2 - 16
)

var (
	ErrPubKeyIDAlreadyExists = errors.New("leaf with given pub key ID already exists")
	ErrInvalidAccountsLength = errors.New("invalid accounts length")
)

type AccountTree struct {
	storage    *Storage
	merkleTree *StoredMerkleTree
}

func NewAccountTree(storage *Storage) *AccountTree {
	return &AccountTree{
		storage:    storage,
		merkleTree: NewStoredMerkleTree("account", storage),
	}
}

func (s *AccountTree) Root() (*common.Hash, error) {
	return s.merkleTree.Root()
}

func (s *AccountTree) LeafNode(pubKeyID uint32) (*models.MerkleTreeNode, error) {
	return s.merkleTree.Get(models.MerklePath{
		Path:  pubKeyID,
		Depth: AccountTreeDepth,
	})
}

func (s *AccountTree) Leaf(pubKeyID uint32) (*models.AccountLeaf, error) {
	leaf, err := s.storage.GetAccountLeaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	return leaf, nil
}

func (s *AccountTree) SetSingle(leaf *models.AccountLeaf) error {
	if leaf.PubKeyID >= leftSubtreeMaxValue {
		return errors.Errorf("invalid pubKeyID value: %d", leaf.PubKeyID)
	}

	tx, storage, err := s.storage.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = NewAccountTree(storage).unsafeSet(leaf)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AccountTree) SetBatch(leaves []models.AccountLeaf) error {
	if len(leaves) != 16 {
		return ErrInvalidAccountsLength
	}

	tx, storage, err := s.storage.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	accountTree := NewAccountTree(storage)

	for i := range leaves {
		if leaves[i].PubKeyID < accountBatchOffset || leaves[i].PubKeyID >= rightSubtreeMaxValue {
			return errors.Errorf("invalid pubKeyID value: %d", leaves[i].PubKeyID)
		}
		_, err = accountTree.unsafeSet(&leaves[i])
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
	err := s.storage.AddAccountLeafIfNotExists(leaf)
	if err == bh.ErrKeyExists {
		return nil, ErrPubKeyIDAlreadyExists
	}
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

func (s *AccountTree) getMerkleTreeNodeByPath(path *models.MerklePath) (*models.MerkleTreeNode, error) {
	node, err := s.merkleTree.Get(*path)
	if err != nil {
		return nil, err
	}
	return node, nil
}
