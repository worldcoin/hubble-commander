package models

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"
)

type testKeyList [][]byte

func TestEncodeKeyList(t *testing.T) {
	keyList := make(testKeyList, 1)
	keyList[0] = append([]byte("_bhIndex:FlatStateLeaf:PubKeyID"), []byte{0, 0, 0, 4}...)
	encoded, err := EncodeKeyList(keyList)
	require.NoError(t, err)

	var decoded testKeyList
	rValue, ok := isKeyListPtrType(&decoded)
	require.True(t, ok)

	err = DecodeKeyList(encoded, *rValue)
	require.NoError(t, err)
	require.Len(t, decoded, len(keyList))
	require.Len(t, decoded[0], len(keyList[0]))
	require.Equal(t, keyList[0], decoded[0])
}

func TestEncodeKeyList_DifferentItemsLength(t *testing.T) {
	keyList := make(testKeyList, 2)
	keyList[0] = append([]byte("_bhIndex:FlatStateLeaf:PubKeyID"), []byte{0, 0, 0, 4}...)
	keyList[1] = append([]byte("_bhIndex:dummy"), []byte{0, 0, 0, 1}...)
	require.Panics(t, func() {
		_, _ = EncodeKeyList(keyList)
	})
}

func TestEncodeKeyList_EmptyKeyList(t *testing.T) {
	keyList := make(testKeyList, 0)
	encoded, err := EncodeKeyList(keyList)
	require.NoError(t, err)
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(encoded[0:4]))
	require.Equal(t, uint32(0), binary.BigEndian.Uint32(encoded[4:8]))

	var decoded testKeyList
	rValue, ok := isKeyListPtrType(&decoded)
	require.True(t, ok)

	err = DecodeKeyList(encoded, *rValue)
	require.NoError(t, err)
	require.Equal(t, keyList, decoded)
}

func TestDecodeKeyList_InvalidLength(t *testing.T) {
	encoded := []byte{0, 0, 0, 1}
	var decoded testKeyList
	rValue, ok := isKeyListPtrType(&decoded)
	require.True(t, ok)

	err := DecodeKeyList(encoded, *rValue)
	require.Error(t, err)
}

func TestDecodeKeyList_InvalidItemLength(t *testing.T) {
	encoded := []byte{0, 0, 0, 1, 0, 0, 0, 1, 0, 2}
	var decoded testKeyList
	rValue, ok := isKeyListPtrType(&decoded)
	require.True(t, ok)

	err := DecodeKeyList(encoded, *rValue)
	require.Error(t, err)
}
