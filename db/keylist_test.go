package db

import (
	"encoding/binary"
	"testing"

	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v4"
)

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
	require.EqualValues(t, len(keyList), binary.BigEndian.Uint32(encoded[0:4]))
	require.EqualValues(t, len(keyList[0]), binary.BigEndian.Uint32(encoded[4:8]))

	var decoded bh.KeyList
	err = DecodeKeyList(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, keyList, decoded)
}

func TestEncodeKeyList_ReturnsErrorWhenItemsHaveInconsistentLengths(t *testing.T) {
	keyList := make(bh.KeyList, 5)
	for i := 0; i < len(keyList)-1; i++ {
		keyList[i] = make([]byte, 4)
		binary.BigEndian.PutUint32(keyList[i], uint32(i))
	}
	keyList[4] = make([]byte, 8)
	binary.BigEndian.PutUint64(keyList[4], uint64(4))

	encoded, err := EncodeKeyList(&keyList)
	require.Nil(t, encoded)
	require.ErrorIs(t, err, ErrInconsistentItemsLength)
}

func TestDecodeKeyList_ReturnsErrorWhenKeyListHasInvalidDataLength(t *testing.T) {
	data := make([]byte, 13)
	data[0] = 2
	data[1] = 5

	var decoded bh.KeyList
	err := DecodeKeyList(data, &decoded)
	require.ErrorIs(t, err, ErrInvalidKeyListLength)
}

func TestKeyListMetadata_Bytes(t *testing.T) {
	metadata := &KeyListMetadata{
		ListLen: 5,
		ItemLen: 10,
	}

	bytes := metadata.Bytes()

	var decoded KeyListMetadata
	err := decoded.SetBytes(bytes)
	require.NoError(t, err)
	require.Equal(t, *metadata, decoded)
}

func TestIndexKeyPrefix(t *testing.T) {
	prefix := IndexKeyPrefix(stored.TxReceiptName, "CommitmentID")
	require.Equal(t, []byte("_bhIndex:TxReceipt:CommitmentID:"), prefix)
}

func TestIndexKey(t *testing.T) {
	value := stored.EncodeUint32(1)
	prefix := IndexKey(stored.TxReceiptName, "CommitmentID", value)
	require.Equal(t, append([]byte("_bhIndex:TxReceipt:CommitmentID:"), value...), prefix)
}
