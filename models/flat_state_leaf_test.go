package models

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestFlatStateLeaf_Bytes(t *testing.T) {
	leaf := FlatStateLeaf{
		StateID:  1,
		DataHash: utils.RandomHash(),
		PubKeyID: 2,
		TokenID:  MakeUint256(3),
		Balance:  MakeUint256(4),
		Nonce:    MakeUint256(5),
	}

	var decodedLeaf FlatStateLeaf
	_ = decodedLeaf.SetBytes(leaf.Bytes())
	require.Equal(t, leaf, decodedLeaf)
}
