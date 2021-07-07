package merkletree

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree_ZeroLeaves(t *testing.T) {
	tree, err := NewMerkleTree([]common.Hash{})
	require.Nil(t, tree)
	require.ErrorIs(t, err, ErrEmptyLeaves)
}

func TestNewMerkleTree_OneNode(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(hash, GetZeroHash(0)), tree.Root())
	require.Equal(t, uint8(2), tree.Depth())
}

func TestNewMerkleTree_TwoNodes(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(left, right), tree.Root())
	require.Equal(t, uint8(2), tree.Depth())
}

func TestNewMerkleTree_ThreeNodes(t *testing.T) {
	leaf1 := utils.RandomHash()
	leaf2 := utils.RandomHash()
	leaf3 := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{leaf1, leaf2, leaf3})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(utils.HashTwo(leaf1, leaf2), utils.HashTwo(leaf3, GetZeroHash(0))), tree.Root())
	require.Equal(t, uint8(3), tree.Depth())
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

	require.Equal(t, models.Witness{right}, tree.GetWitness(0))
	require.Equal(t, models.Witness{left}, tree.GetWitness(1))
}

func TestMerkleTree_GetWitness_ThreeNodes(t *testing.T) {
	leaf1 := utils.RandomHash()
	leaf2 := utils.RandomHash()
	leaf3 := utils.RandomHash()

	h12 := utils.HashTwo(leaf1, leaf2)
	h30 := utils.HashTwo(leaf3, GetZeroHash(0))

	tree, err := NewMerkleTree([]common.Hash{leaf1, leaf2, leaf3})
	require.NoError(t, err)

	require.Equal(t, models.Witness{leaf2, h30}, tree.GetWitness(0))
	require.Equal(t, models.Witness{leaf1, h30}, tree.GetWitness(1))
	require.Equal(t, models.Witness{GetZeroHash(0), h12}, tree.GetWitness(2))
}
