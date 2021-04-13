package dto

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestCreate2Transfer_JSONMarshaling(t *testing.T) {
	transfer := Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		ToPubkeyID:  ref.Uint32(3),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   utils.RandomBytes(12),
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	var unmarshalled Create2Transfer
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, transfer, unmarshalled)
}

func TestCreate2Transfer_MarshalJSON(t *testing.T) {
	sig, err := hex.DecodeString("deadbeef")
	require.NoError(t, err)

	transfer := Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		ToPubkeyID:  ref.Uint32(3),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   sig,
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	// nolint:goconst,lll
	expected := `{"Type":3,"FromStateID":1,"ToStateID":2,"ToPubkeyID":3,"Amount":"50","Fee":"10","Nonce":"0","Signature":"0xdeadbeef"}`
	require.Equal(t, expected, string(data))
}
