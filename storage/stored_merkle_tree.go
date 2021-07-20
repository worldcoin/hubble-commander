package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

type StoredMerkleTree struct {
	storage *Storage
	prefix  string
}

func NewStoredMerkleTree(prefix string, storage *Storage) *StoredMerkleTree {
	return &StoredMerkleTree{
		storage: storage,
		prefix:  prefix,
	}
}

func (s *StoredMerkleTree) keyFor(path models.MerklePath) models.MerklePathWithPrefix {
	return models.MerklePathWithPrefix{Prefix: s.prefix, Path: path}
}

func (s *StoredMerkleTree) Get(path models.MerklePath) (*models.StateNode, error) {
	node := models.StateNode{MerklePath: path}
	err := s.storage.Badger.Get(s.keyFor(path), &node)
	if err == bh.ErrNotFound {
		return newZeroStateNode(&path), nil
	}
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *StoredMerkleTree) Root() (*common.Hash, error) {
	node, err := s.Get(models.MerklePath{Path: 0, Depth: 0})
	if err != nil {
		return nil, err
	}

	return &node.DataHash, nil
}

func (s *StoredMerkleTree) SetSingleNode(node *models.StateNode) error {
	return s.storage.Badger.Upsert(s.keyFor(node.MerklePath), *node)
}

// SetNode sets node hash and update all nodes leading to root. Returns new root hash and the insertion witness.
func (s *StoredMerkleTree) SetNode(path *models.MerklePath, hash *common.Hash) (*common.Hash, models.Witness, error) {
	currentPath := path
	currentHash := *hash
	witness := make(models.Witness, 0, path.Depth)

	for currentPath.Depth != 0 {
		sibling, err := currentPath.Sibling()
		if err != nil {
			return nil, nil, err
		}

		siblingNode, err := s.storage.GetStateNodeByPath(sibling)
		if err != nil {
			return nil, nil, err
		}
		witness = append(witness, siblingNode.DataHash)

		err = s.SetSingleNode(&models.StateNode{
			MerklePath: *currentPath,
			DataHash:   currentHash,
		})
		if err != nil {
			return nil, nil, err
		}
		currentHash = calculateParentHash(&currentHash, currentPath, siblingNode.DataHash)

		currentPath, err = currentPath.Parent()
		if err != nil {
			return nil, nil, err
		}
	}

	rootPath := models.MerklePath{Depth: 0, Path: 0}
	err := s.storage.UpsertStateNode(&models.StateNode{
		MerklePath: rootPath,
		DataHash:   currentHash,
	})
	if err != nil {
		return nil, nil, err
	}

	return &currentHash, witness, nil
}

func calculateParentHash(
	currentHash *common.Hash,
	currentPath *models.MerklePath,
	witnessHash common.Hash,
) common.Hash {
	if currentPath.IsLeftNode() {
		return utils.HashTwo(*currentHash, witnessHash)
	} else {
		return utils.HashTwo(witnessHash, *currentHash)
	}
}
