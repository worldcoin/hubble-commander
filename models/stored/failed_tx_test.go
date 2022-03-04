package stored

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func TestNewFailedTxIndex(t *testing.T) {
	fromStateID := uint32(1)
	nonce := models.NewUint256(11)
	failedTxIndex := NewFailedTxIndex(fromStateID, nonce)

	var decodedFromStateID uint32
	err := DecodeUint32(failedTxIndex[:4], &decodedFromStateID)
	require.NoError(t, err)
	require.Equal(t, fromStateID, decodedFromStateID)

	decodedNonce := new(models.Uint256)
	decodedNonce.SetBytes(failedTxIndex[4:])
	require.Equal(t, nonce, decodedNonce)
}
