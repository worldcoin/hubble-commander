package storage

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestZeroHash_Root(t *testing.T) {
	require.Equal(t, common.HexToHash("0xcf277fb80a82478460e8988570b718f1e083ceb76f7e271a1a1497e5975f53ae"), GetZeroHash(leafDepth))
}

func TestZeroHash_RootChild(t *testing.T) {
	require.Equal(t, common.HexToHash("0x78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62"), GetZeroHash(31))
}

func TestZeroHash_Panic(t *testing.T) {
	require.Panics(t, func() { GetZeroHash(33) })
}
