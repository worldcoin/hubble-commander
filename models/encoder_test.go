package models

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_EncodePointer(t *testing.T) {
	leaf := &FlatStateLeaf{
		StateID:  5,
		DataHash: common.Hash{1, 2, 3},
		PubKeyID: 3,
		TokenID:  MakeUint256(1),
		Balance:  MakeUint256(9999888),
		Nonce:    MakeUint256(0),
	}

	bytes := encodePointer(136, leaf)
	require.EqualValues(t, 1, bytes[0])
	require.Equal(t, leaf.Bytes(), bytes[1:])
}

func Test_EncodePointer_NilValue(t *testing.T) {
	var leaf *FlatStateLeaf
	bytes := encodePointer(136, leaf)
	require.EqualValues(t, bytes[0], 0)
}

func Test_EncodeHashPointer(t *testing.T) {
	hash := &common.Hash{1, 2, 3, 4}
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := DecodeHashPointer(bytes)
	require.Equal(t, *hash, *decodedValue)
}

func Test_EncodeHashPointer_NilValue(t *testing.T) {
	var hash *common.Hash
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := DecodeHashPointer(bytes)
	require.Nil(t, decodedValue)
}

func Test_EncodeUint32Pointer(t *testing.T) {
	value := uint32(32)
	bytes := encodeUint32Pointer(&value)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := decodeUint32Pointer(bytes)
	require.Equal(t, value, *decodedValue)
}

func Test_EncodeUint32Pointer_NilValue(t *testing.T) {
	var value *uint32
	bytes := encodeUint32Pointer(value)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := decodeUint32Pointer(bytes)
	require.Nil(t, decodedValue)
}
