package badger

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestDataHash_ByteEncoding(t *testing.T) {
	node := models.StateNode{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}

	var decodedNode models.StateNode
	encodedDataHash, _ := EncodeDataHash(&node)
	_ = DecodeDataHash(encodedDataHash, &decodedNode)
	require.Equal(t, node, decodedNode)
}

func TestUint32_ByteEncoding(t *testing.T) {
	number := uint32(173)

	var decodedNumber uint32
	encodedDataHash, _ := EncodeUint32(&number)
	_ = DecodeUint32(encodedDataHash, &decodedNumber)
	require.Equal(t, number, decodedNumber)
}
