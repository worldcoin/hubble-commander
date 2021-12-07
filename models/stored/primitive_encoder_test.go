package stored

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestEncodeHash(t *testing.T) {
	dataHash := common.BytesToHash([]byte{1, 2, 3, 4, 5})

	var decodedHash common.Hash
	encodedDataHash, _ := EncodeHash(&dataHash)
	_ = DecodeDataHash(encodedDataHash, &decodedHash)
	require.Equal(t, dataHash, decodedHash)
}

func TestEncodeUint32(t *testing.T) {
	number := uint32(173)

	var decodedNumber uint32
	encodedDataHash := EncodeUint32(number)
	_ = DecodeUint32(encodedDataHash, &decodedNumber)
	require.Equal(t, number, decodedNumber)
}

func TestEncodeUint64(t *testing.T) {
	value := uint64(123456789)

	encoded := EncodeUint64(value)

	var decoded uint64
	err := DecodeUint64(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}

func TestEncodeString(t *testing.T) {
	encoded := EncodeString(testMessage)

	var decoded string
	err := DecodeString(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, testMessage, decoded)
}
