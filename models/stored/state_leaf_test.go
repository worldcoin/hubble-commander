package stored

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestStateLeaf_Bytes(t *testing.T) {
	leaf := FlatStateLeaf{
		StateID:  1,
		DataHash: utils.RandomHash(),
		PubKeyID: 2,
		TokenID:  models.MakeUint256(3),
		Balance:  models.MakeUint256(4),
		Nonce:    models.MakeUint256(5),
	}

	var decodedLeaf FlatStateLeaf
	_ = decodedLeaf.SetBytes(leaf.Bytes())
	require.Equal(t, leaf, decodedLeaf)
}
