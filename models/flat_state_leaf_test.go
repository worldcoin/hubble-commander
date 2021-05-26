package models

import (
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
	leaf := FlatStateLeaf{
		PubKeyID:   26,
		TokenIndex: MakeUint256(25),
	}
	encoded, err := tupleIndex("Tuple", leaf)
	require.NoError(t, err)

	var decoded StateLeafIndex
	err = Decode(encoded, &decoded)
	require.NoError(t, err)
	require.Equal(t, StateLeafIndex{
		PubKeyID:   leaf.PubKeyID,
		TokenIndex: leaf.TokenIndex,
	}, decoded)
}
