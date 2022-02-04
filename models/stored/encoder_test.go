package stored

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const testMessage = "some message"

func TestEncodeHashPointer(t *testing.T) {
	hash := &common.Hash{1, 2, 3, 4}
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := decodeHashPointer(bytes)
	require.Equal(t, *hash, *decodedValue)
}

func TestEncodeHashPointer_NilValue(t *testing.T) {
	var hash *common.Hash
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := decodeHashPointer(bytes)
	require.Nil(t, decodedValue)
}

func TestEncodeUint32Pointer(t *testing.T) {
	value := uint32(32)
	bytes := EncodeUint32Pointer(&value)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := decodeUint32Pointer(bytes)
	require.Equal(t, value, *decodedValue)
}

func TestEncodeUint32Pointer_NilValue(t *testing.T) {
	var value *uint32
	bytes := EncodeUint32Pointer(value)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := decodeUint32Pointer(bytes)
	require.Nil(t, decodedValue)
}

func TestEncodeTimestampPointer(t *testing.T) {
	timestamp := models.NewTimestamp(time.Unix(10, 0).UTC())
	bytes := encodeTimestampPointer(timestamp)
	require.EqualValues(t, 1, bytes[0])

	decodedTimestamp, err := decodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Equal(t, *timestamp, *decodedTimestamp)
}

func TestEncodeTimestampPointer_NilValue(t *testing.T) {
	var timestamp *models.Timestamp
	bytes := encodeTimestampPointer(timestamp)
	require.EqualValues(t, 0, bytes[0])

	decodedTimestamp, err := decodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Nil(t, decodedTimestamp)
}
