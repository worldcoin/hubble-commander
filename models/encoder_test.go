package models

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

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

func Test_encodeStringPointer(t *testing.T) {
	value := "some string"
	bytes := encodeStringPointer(&value)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := decodeStringPointer(bytes)
	require.Equal(t, value, *decodedValue)
}

func Test_EncodeStringPointer_NilValue(t *testing.T) {
	var value *string
	bytes := encodeStringPointer(value)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := decodeStringPointer(bytes)
	require.Nil(t, decodedValue)
}

func Test_EncodeTimestampPointer(t *testing.T) {
	timestamp := NewTimestamp(time.Unix(10, 0).UTC())
	bytes := encodeTimestampPointer(timestamp)
	require.EqualValues(t, 1, bytes[0])

	decodedTimestamp, err := decodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Equal(t, *timestamp, *decodedTimestamp)
}

func Test_EncodeTimestampPointer_NilValue(t *testing.T) {
	var timestamp *Timestamp
	bytes := encodeTimestampPointer(timestamp)
	require.EqualValues(t, 0, bytes[0])

	decodedTimestamp, err := decodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Nil(t, decodedTimestamp)
}
