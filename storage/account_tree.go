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

	batchSize            = 1 << 4
	accountBatchOffset   = 1 << 31
	leftSubtreeMaxValue  = accountBatchOffset - 2
	rightSubtreeMaxValue = accountBatchOffset*2 - 18
)

var ErrInvalidAccountsLength = errors.New("invalid accounts length")

type AccountTree struct {
	storageBase *StorageBase
	merkleTree  *StoredMerkleTree
}

func NewAccountTree(storageBase *StorageBase) *AccountTree {
	return &AccountTree{
		storageBase: storageBase,
		merkleTree:  NewStoredMerkleTree("account", storageBase.Badger),
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
	var leaf models.AccountLeaf
	err := s.storageBase.Badger.Get(pubKeyID, &leaf)
	if err == bh.ErrNotFound {
		return nil, NewNotFoundError("account leaf")
	}
	if err != nil {
		return nil, err
	}
	return &leaf, nil
}

func (s *AccountTree) Leaves(publicKey *models.PublicKey) ([]models.AccountLeaf, error) {
	return s.getAccountLeaves(publicKey)
}

func (s *AccountTree) SetSingle(leaf *models.AccountLeaf) error {
	if leaf.PubKeyID > leftSubtreeMaxValue {
		return NewInvalidPubKeyIDError(leaf.PubKeyID)
	}

	tx, storage, err := s.storageBase.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	_, err = NewAccountTree(storage).unsafeSet(leaf)
	if err == bh.ErrKeyExists {
		return NewAccountAlreadyExistsError(leaf)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *AccountTree) SetBatch(leaves []models.AccountLeaf) error {
	if len(leaves) != batchSize {
		return ErrInvalidAccountsLength
	}

	tx, storage, err := s.storageBase.BeginTransaction(TxOptions{Badger: true})
	if err != nil {
		return err
	}
	defer tx.Rollback(&err)

	accountTree := NewAccountTree(storage)

	for i := range leaves {
		if leaves[i].PubKeyID < accountBatchOffset || leaves[i].PubKeyID > rightSubtreeMaxValue {
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
	err := s.storageBase.Badger.Insert(leaf.PubKeyID, *leaf)
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
	return s.merkleTree.Get(*path)
}
