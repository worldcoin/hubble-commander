package merkletree

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_GetZeroHash_Root(t *testing.T) {
	require.Equal(
		t,
		common.HexToHash("0xcf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae"),
		GetZeroHash(MaxDepth),
	)
}

func Test_GetZeroHash_RootChild(t *testing.T) {
	require.Equal(
		t,
		common.HexToHash("0x78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62"),
		GetZeroHash(31),
	)
}

func Test_GetZeroHash_Panic(t *testing.T) {
	require.Panics(t, func() { GetZeroHash(MaxDepth + 1) })
}

func Test_NewMerkleTree_ZeroLeaves(t *testing.T) {
	tree, err := NewMerkleTree([]common.Hash{})
	require.Nil(t, tree)
	require.ErrorIs(t, err, ErrEmptyLeaves)
}

func Test_NewMerkleTree_OneLeaf(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(hash, GetZeroHash(0)), tree.Root())
	require.Equal(t, uint8(2), tree.Depth())
}

func Test_NewMerkleTree_TwoLeaves(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(left, right), tree.Root())
	require.Equal(t, uint8(2), tree.Depth())
}

func Test_NewMerkleTree_ThreeLeaves(t *testing.T) {
	leaf1 := utils.RandomHash()
	leaf2 := utils.RandomHash()
	leaf3 := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{leaf1, leaf2, leaf3})
	require.NoError(t, err)

	require.Equal(t, utils.HashTwo(utils.HashTwo(leaf1, leaf2), utils.HashTwo(leaf3, GetZeroHash(0))), tree.Root())
	require.Equal(t, uint8(3), tree.Depth())
}

func TestMerkleTree_GetWitness_OneLeaf(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	require.Equal(t, models.Witness{GetZeroHash(0)}, tree.GetWitness(0))
}

func TestMerkleTree_GetWitness_TwoLeaves(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, models.Witness{right}, tree.GetWitness(0))
	require.Equal(t, models.Witness{left}, tree.GetWitness(1))
}

func TestMerkleTree_GetWitness_ThreeLeaves(t *testing.T) {
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
