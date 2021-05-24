package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	bh "github.com/timshannon/badgerhold/v3"
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
	encoded, err := PubKeyIDIndex("PubKeyID", leaf)
	require.NoError(t, err)

	var decoded uint32
	err = bh.DefaultDecode(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, leaf.PubKeyID, decoded)
}
