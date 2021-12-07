package merkletree

import (
	stdErr "errors"
	"fmt"
	"math"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

const MaxDepth = 32

var (
	zeroHashes [MaxDepth + 1]common.Hash

	ErrEmptyLeaves   = stdErr.New("leaves cannot be empty")
	ErrTooManyLeaves = stdErr.New("merkle tree too large")
)

func init() {
	zeroHashes[0] = consts.ZeroHash
	for i := 1; i <= MaxDepth; i++ {
		zeroHashes[i] = utils.HashTwo(zeroHashes[i-1], zeroHashes[i-1])
	}
}

func GetZeroHash(level uint8) common.Hash {
	if level > MaxDepth {
		panic(fmt.Sprintf("level > %d", MaxDepth))
	}

	return zeroHashes[level]
}

type MerkleTree struct {
	// Nodes are stored in a flat array which starts with the root,
	// then the 2 nodes of the first level left to right,
	// then the next 4 nodes, and so on finishing with the array of leaves.
	//
	// The valid lengths for the nodes array should be in form of 2^n - 1 (i.e. 1, 3, 7 elements and so on).
	nodes []common.Hash

	// Depth of the tree, equal to the number of layers in the tree.
	depth uint8
}

func NewMerkleTree(leaves []common.Hash) (*MerkleTree, error) {
	if len(leaves) == 0 {
		return nil, errors.WithStack(ErrEmptyLeaves)
	}
	depth := getRequiredTreeHeight(int32(len(leaves)))

	if depth > MaxDepth {
		return nil, errors.WithStack(ErrTooManyLeaves)
	}

	arraySize := (1 << depth) - 1
	tree := &MerkleTree{
		nodes: make([]common.Hash, arraySize),
		depth: depth,
	}

	// Set the known leaves.
	for i := range leaves {
		tree.setNode(models.MerklePath{Depth: depth - 1, Path: uint32(i)}, leaves[i])
	}

	// Set the rest of the leaves on the lowest level to "zero hash".
	for i := len(leaves); uint32(i) < getNodeCountAtDepth(depth-1); i++ {
		tree.setNode(models.MerklePath{Depth: depth - 1, Path: uint32(i)}, GetZeroHash(0))
	}

	// Populate the rest of the levels with hashes of their children.
	for level := int(depth) - 2; level >= 0; level-- {
		nodeCount := getNodeCountAtDepth(uint8(level))
		for i := 0; uint32(i) < nodeCount; i++ {
			path := models.MerklePath{Depth: uint8(level), Path: uint32(i)}
			leftPath, err := path.Child(false)
			if err != nil {
				return nil, err
			}
			leftHash := tree.GetNode(*leftPath)

			rightPath, err := path.Child(true)
			if err != nil {
				return nil, err
			}

			rightHash := tree.GetNode(*rightPath)

			tree.setNode(path, utils.HashTwo(leftHash, rightHash))
		}
	}

	return tree, nil
}

func (m *MerkleTree) setNode(path models.MerklePath, value common.Hash) {
	m.nodes[getNodeIndex(path)] = value
}

func (m *MerkleTree) Depth() uint8 {
	return m.depth
}

func (m *MerkleTree) GetNode(path models.MerklePath) common.Hash {
	return m.nodes[getNodeIndex(path)]
}

func (m *MerkleTree) Root() common.Hash {
	return m.nodes[0]
}

func (m *MerkleTree) GetWitness(leafIndex uint32) models.Witness {
	leafPath := models.MerklePath{Depth: m.depth - 1, Path: leafIndex}

	witness := make([]common.Hash, m.depth-1)
	for leafPath.Depth > 0 {
		sibling, err := leafPath.Sibling()
		if err != nil {
			panic(err) // Can not fail
		}

		witness[m.depth-leafPath.Depth-1] = m.GetNode(*sibling)

		parent, err := leafPath.Parent()
		if err != nil {
			panic(err) // Can not fail
		}
		leafPath = *parent
	}

	return witness
}

func getRequiredTreeHeight(leafCount int32) uint8 {
	if leafCount == 1 {
		return 2
	}
	return uint8(math.Ceil(math.Log2(float64(leafCount)))) + 1
}

func getNodeIndex(path models.MerklePath) int {
	return int((1 << path.Depth) + path.Path - 1)
}

func getNodeCountAtDepth(depth uint8) uint32 {
	return 1 << depth
}
