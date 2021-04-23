package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestCreate2Transfer_JSONMarshaling(t *testing.T) {
	transfer := Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: &models.PublicKey{1, 2, 3},
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
		ToPublicKey: &models.PublicKey{1, 2, 3},
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{4, 5, 6},
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	expected := fmt.Sprintf(
		`{"Type":"CREATE2TRANSFER","FromStateID":1,"ToPublicKey":"%s","Amount":"50","Fee":"10","Nonce":"0","Signature":"%s"}`,
		"0x010203"+strings.Repeat("0", models.PublicKeyLength*2-6),
		"0x040506"+strings.Repeat("0", models.SignatureLength*2-6),
	)
	require.Equal(t, expected, string(data))
}
