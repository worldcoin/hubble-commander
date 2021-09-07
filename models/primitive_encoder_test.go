package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEncodeDataHash(t *testing.T) {
	node := MerkleTreeNode{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}

	var decodedHash common.Hash
	encodedDataHash, _ := EncodeDataHash(&node.DataHash)
	_ = DecodeDataHash(encodedDataHash, &decodedHash)
	require.Equal(t, node.DataHash, decodedHash)
}

func TestEncodeUint32(t *testing.T) {
	number := uint32(173)

	var decodedNumber uint32
	encodedDataHash, _ := EncodeUint32(&number)
	_ = DecodeUint32(encodedDataHash, &decodedNumber)
	require.Equal(t, number, decodedNumber)
}

func TestEncodeUint64(t *testing.T) {
	value := uint64(123456789)

	encoded, err := EncodeUint64(&value)
	require.NoError(t, err)

	var decoded uint64
	err = DecodeUint64(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}

func TestEncodeString(t *testing.T) {
	encoded, err := EncodeString(ref.String(testMessage))
	require.NoError(t, err)

	var decoded string
	err = DecodeString(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, testMessage, decoded)
}
