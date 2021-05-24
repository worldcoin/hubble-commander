package models

import (
	"encoding/binary"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestFlatStateLeaf_ByteEncoding(t *testing.T) {
	leaf := FlatStateLeaf{
		StateID:    1,
		DataHash:   utils.RandomHash(),
		PubKeyID:   2,
		TokenIndex: MakeUint256(3),
		Balance:    MakeUint256(4),
		Nonce:      MakeUint256(5),
	}

	var decodedLeaf FlatStateLeaf
	_ = decodedLeaf.SetBytes(leaf.Bytes())
	require.Equal(t, leaf, decodedLeaf)
}

func Test_pubKeyIDIndex(t *testing.T) {
	leaf := FlatStateLeaf{PubKeyID: 26}
	encoded, err := pubKeyIDIndex("PubKeyID", leaf)
	require.NoError(t, err)

	decoded := binary.BigEndian.Uint32(encoded)
	require.Equal(t, leaf.PubKeyID, decoded)
}
