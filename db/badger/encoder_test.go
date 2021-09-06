package badger

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func TestDecodeKey(t *testing.T) {
	prefix := []byte("bh_prefix")
	value := uint64(123456789)

	encoded, err := models.EncodeUint64(&value)
	require.NoError(t, err)

	var decoded uint64
	err = DecodeKey(append(prefix, encoded...), &decoded, prefix)
	require.NoError(t, err)
	require.Equal(t, value, decoded)
}
