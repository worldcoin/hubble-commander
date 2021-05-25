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

func TestUint64_ByteEncoding(t *testing.T) {
	value := uint64(123456789)

	encoded, err := EncodeUint64(&value)
	require.NoError(t, err)

	var decoded uint64
	err = DecodeUint64(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}

func TestDecodeKey(t *testing.T) {
	prefix := []byte("bh_prefix")
	value := uint64(123456789)

	encoded, err := EncodeUint64(&value)
	require.NoError(t, err)

	var decoded uint64
	err = DecodeKey(append(prefix, encoded...), &decoded, prefix)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}
