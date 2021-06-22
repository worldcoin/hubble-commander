package merkle_tree

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree_OnlyRoot(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	require.Equal(t, hash, tree.Root())
}

func TestNewMerkleTree_TwoNodes(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(left, right), tree.Root())
}

func TestNewMerkleTree_ThreeNodes(t *testing.T) {
	// TODO
}

func TestMerkleTree_GetWitness_OnlyRoot(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	witness := tree.GetWitness(0)
	require.Len(t, witness, 0)
}

func TestMerkleTree_GetWitness_TwoNodes(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, Witness{right}, tree.GetWitness(0))
	require.Equal(t, Witness{left}, tree.GetWitness(1))
}

func TestMerkleTree_GetWitness_ThreeNodes(t *testing.T) {
	// TODO
}
