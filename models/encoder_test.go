package models

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const testMessage = "some message"

func TestEncodeHashPointer(t *testing.T) {
	hash := &common.Hash{1, 2, 3, 4}
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := DecodeHashPointer(bytes)
	require.Equal(t, *hash, *decodedValue)
}

func TestEncodeHashPointer_NilValue(t *testing.T) {
	var hash *common.Hash
	bytes := EncodeHashPointer(hash)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := DecodeHashPointer(bytes)
	require.Nil(t, decodedValue)
}

func TestEncodeUint32Pointer(t *testing.T) {
	value := uint32(32)
	bytes := EncodeUint32Pointer(&value)
	require.EqualValues(t, 1, bytes[0])

	decodedValue := DecodeUint32Pointer(bytes)
	require.Equal(t, value, *decodedValue)
}

func TestEncodeUint32Pointer_NilValue(t *testing.T) {
	var value *uint32
	bytes := EncodeUint32Pointer(value)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := DecodeUint32Pointer(bytes)
	require.Nil(t, decodedValue)
}

func TestEncodeStringPointer(t *testing.T) {
	bytes := EncodeStringPointer(ref.String(testMessage))
	require.EqualValues(t, 1, bytes[0])

	decodedValue := DecodeStringPointer(bytes)
	require.Equal(t, testMessage, *decodedValue)
}

func TestEncodeStringPointer_NilValue(t *testing.T) {
	var value *string
	bytes := EncodeStringPointer(value)
	require.EqualValues(t, 0, bytes[0])

	decodedValue := DecodeStringPointer(bytes)
	require.Nil(t, decodedValue)
}

func TestEncodeTimestampPointer(t *testing.T) {
	timestamp := NewTimestamp(time.Unix(10, 0).UTC())
	bytes := EncodeTimestampPointer(timestamp)
	require.EqualValues(t, 1, bytes[0])

	decodedTimestamp, err := DecodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Equal(t, *timestamp, *decodedTimestamp)
}

func TestEncodeTimestampPointer_NilValue(t *testing.T) {
	var timestamp *Timestamp
	bytes := EncodeTimestampPointer(timestamp)
	require.EqualValues(t, 0, bytes[0])

	decodedTimestamp, err := DecodeTimestampPointer(bytes)
	require.NoError(t, err)
	require.Nil(t, decodedTimestamp)
}

func TestEncodeCommitmentIDPointer(t *testing.T) {
	id := &CommitmentID{
		BatchID:      MakeUint256(5),
		IndexInBatch: 2,
	}
	bytes := EncodeCommitmentIDPointer(id)
	require.EqualValues(t, 1, bytes[0])

	decodedID, err := DecodeCommitmentIDPointer(bytes)
	require.NoError(t, err)
	require.Equal(t, *id, *decodedID)
}

func TestEncodeCommitmentIDPointer_NilValue(t *testing.T) {
	var id *CommitmentID
	bytes := EncodeCommitmentIDPointer(id)
	require.EqualValues(t, 0, bytes[0])

	decodedID, err := DecodeCommitmentIDPointer(bytes)
	require.NoError(t, err)
	require.Nil(t, decodedID)
}
