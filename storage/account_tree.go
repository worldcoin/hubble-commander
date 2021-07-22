package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const AccountTreeDepth = merkletree.MaxDepth

var ErrPubKeyIDAlreadyExists = errors.New("leaf with given pub key ID already exists")

type AccountTree struct {
	storage    *Storage
	merkleTree *StoredMerkleTree
}

func NewAccountTree(storage *Storage) *AccountTree {
	return &AccountTree{
		storage:    storage,
		merkleTree: NewStoredMerkleTree("state", storage),
	}
}

func (s *AccountTree) Root() (*common.Hash, error) {
	return s.merkleTree.Root()
}

func (s *AccountTree) LeafNode(pubKeyID uint32) (*models.MerkleTreeNode, error) {
	return s.merkleTree.Get(models.MerklePath{
		Path:  pubKeyID,
		Depth: StateTreeDepth,
	})
}

func (s *AccountTree) Leaf(pubKeyID uint32) (*models.AccountLeaf, error) {
	leaf, err := s.storage.GetAccountLeaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	return leaf, nil
}

// Set returns a witness containing 32 elements for the current set operation
func (s *AccountTree) Set(leaf *models.AccountLeaf) (models.Witness, error) {
	tx, storage, err := s.storage.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(&err)

	witness, err := NewAccountTree(storage).unsafeSet(leaf)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return witness, nil
}

func (s *AccountTree) GetWitness(id uint32) (models.Witness, error) {
	return s.merkleTree.GetWitness(models.MakeMerklePathFromLeafID(id))
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
