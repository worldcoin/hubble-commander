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

func TestTransfer_JSONMarshaling(t *testing.T) {
	transfer := Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   models.NewRandomSignature(),
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	var unmarshalled Transfer
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, transfer, unmarshalled)
}

func TestTransfer_MarshalJSON(t *testing.T) {
	transfer := Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{1, 2, 3},
	}
	data, err := json.Marshal(transfer)
	require.NoError(t, err)

	expected := fmt.Sprintf(
		`{"Type":"TRANSFER","FromStateID":1,"ToStateID":2,"Amount":"50","Fee":"10","Nonce":"0","Signature":"%s"}`,
		"0x010203"+strings.Repeat("0", models.SignatureLength*2-6),
	)
	require.Equal(t, expected, string(data))
}
