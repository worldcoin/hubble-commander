package merkletree

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree_OnlyRoot(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	require.Equal(t, hash, tree.Root())
	require.Equal(t, uint8(1), tree.Depth())
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

	require.Equal(t, utils.HashTwo(utils.HashTwo(leaf1, leaf2), utils.HashTwo(leaf3, GetZeroHash(2))), tree.Root())
	require.Equal(t, uint8(3), tree.Depth())
}

func TestMerkleTree_GetWitness_OnlyRoot(t *testing.T) {
	hash := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{hash})
	require.NoError(t, err)

	witness := tree.GetWitnesses(0)
	require.Len(t, witness, 0)
}

func TestMerkleTree_GetWitness_TwoNodes(t *testing.T) {
	left := utils.RandomHash()
	right := utils.RandomHash()

	tree, err := NewMerkleTree([]common.Hash{left, right})
	require.NoError(t, err)

	require.Equal(t, models.Witnesses{right}, tree.GetWitnesses(0))
	require.Equal(t, models.Witnesses{left}, tree.GetWitnesses(1))
}

func TestMerkleTree_GetWitness_ThreeNodes(t *testing.T) {
	leaf1 := utils.RandomHash()
	leaf2 := utils.RandomHash()
	leaf3 := utils.RandomHash()

	h12 := utils.HashTwo(leaf1, leaf2)
	h30 := utils.HashTwo(leaf3, GetZeroHash(2))

	tree, err := NewMerkleTree([]common.Hash{leaf1, leaf2, leaf3})
	require.NoError(t, err)

	require.Equal(t, models.Witnesses{leaf2, h30}, tree.GetWitnesses(0))
	require.Equal(t, models.Witnesses{leaf1, h30}, tree.GetWitnesses(1))
	require.Equal(t, models.Witnesses{GetZeroHash(2), h12}, tree.GetWitnesses(2))
}
