package dto

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestCreate2Transfer_JSONMarshaling(t *testing.T) {
	transfer := Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: &examplePublicKey,
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   models.NewRandomSignature(),
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	var unmarshalled Create2Transfer
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, transfer, unmarshalled)
}

func TestCreate2Transfer_MarshalJSON(t *testing.T) {
	transfer := Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: &examplePublicKey,
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &exampleSignature,
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	expected := fmt.Sprintf(
		`{"Type":"CREATE2TRANSFER","FromStateID":1,"ToPublicKey":%q,"Amount":"50","Fee":"10","Nonce":"0","Signature":%q}`,
		examplePublicKeyHex,
		exampleSignatureHex,
	)
	require.Equal(t, expected, string(data))
}
