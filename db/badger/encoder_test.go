package badger

import (
	"encoding/binary"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v3"
)

func TestDataHash_ByteEncoding(t *testing.T) {
	node := models.MerkleTreeNode{
		DataHash: common.BytesToHash([]byte{1, 2, 3, 4, 5}),
	}

	var decodedHash common.Hash
	encodedDataHash, _ := EncodeDataHash(&node.DataHash)
	_ = DecodeDataHash(encodedDataHash, &decodedHash)
	require.Equal(t, node.DataHash, decodedHash)
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

func TestEncodeKeyList(t *testing.T) {
	prefix := []byte("bh_prefix")
	keyList := make(bh.KeyList, 5)
	for i := range keyList {
		keyList[i] = make([]byte, len(prefix)+4)
		copy(keyList[i][:len(prefix)], prefix)
		binary.BigEndian.PutUint32(keyList[i][len(prefix):], uint32(i))
	}

	encoded, err := EncodeKeyList(&keyList)
	require.NoError(t, err)
	require.EqualValues(t, len(keyList), encoded[0])
	require.EqualValues(t, len(keyList[0]), encoded[1])

	var decoded bh.KeyList
	err = DecodeKeyList(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, keyList, decoded)
}

func TestEncodeString(t *testing.T) {
	value := "some string"

	encoded, err := EncodeString(&value)
	require.NoError(t, err)

	var decoded string
	err = DecodeString(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}
