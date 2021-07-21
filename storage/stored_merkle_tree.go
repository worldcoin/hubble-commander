package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	bh "github.com/timshannon/badgerhold/v3"
)

type StoredMerkleTree struct {
	storage   *Storage
	namespace string
}

func NewStoredMerkleTree(namespace string, storage *Storage) *StoredMerkleTree {
	return &StoredMerkleTree{
		storage:   storage,
		namespace: namespace,
	}
}

func (s *StoredMerkleTree) keyFor(path models.MerklePath) models.NamespacedMerklePath {
	return models.NamespacedMerklePath{Namespace: s.namespace, Path: path}
}

func (s *StoredMerkleTree) Get(path models.MerklePath) (*models.MerkleTreeNode, error) {
	node := models.MerkleTreeNode{MerklePath: path}
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

func (s *StoredMerkleTree) SetSingleNode(node *models.MerkleTreeNode) error {
	return s.storage.Badger.Upsert(s.keyFor(node.MerklePath), *node)
}

// SetNode sets node hash and update all nodes leading to root. Returns new root hash and the insertion witness.
func (s *StoredMerkleTree) SetNode(path *models.MerklePath, hash common.Hash) (*common.Hash, models.Witness, error) {
	currentPath := path
	currentHash := hash
	witness := make(models.Witness, 0, path.Depth)

	for currentPath.Depth != 0 {
		sibling, err := currentPath.Sibling()
		if err != nil {
			return nil, nil, err
		}

		siblingNode, err := s.Get(*sibling)
		if err != nil {
			return nil, nil, err
		}
		witness = append(witness, siblingNode.DataHash)

		err = s.SetSingleNode(&models.MerkleTreeNode{
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
	err := s.SetSingleNode(&models.MerkleTreeNode{
		MerklePath: rootPath,
		DataHash:   currentHash,
	})
	if err != nil {
		return nil, nil, err
	}

	return &currentHash, witness, nil
}

func (s *StoredMerkleTree) GetWitness(path models.MerklePath) (models.Witness, error) {
	witnessPaths, err := path.GetWitnessPaths()
	if err != nil {
		return nil, err
	}

	witness := make([]common.Hash, 0, len(witnessPaths))
	for i := range witnessPaths {
		var node *models.MerkleTreeNode
		node, err = s.Get(witnessPaths[i])
		if err != nil {
			return nil, err
		}
		witness = append(witness, node.DataHash)
	}

	return witness, nil
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

func newZeroStateNode(path *models.MerklePath) *models.MerkleTreeNode {
	return &models.MerkleTreeNode{
		MerklePath: *path,
		DataHash:   merkletree.GetZeroHash(StateTreeDepth - uint(path.Depth)),
	}
}
