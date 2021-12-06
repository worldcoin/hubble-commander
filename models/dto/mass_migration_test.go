package dto

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestMassMigration_JSONMarshaling(t *testing.T) {
	massMigration := MassMigration{
		FromStateID: ref.Uint32(1),
		SpokeID:     ref.Uint32(5),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   models.NewRandomSignature(),
	}
	data, err := json.Marshal(massMigration)
	require.NoError(t, err)

	var unmarshalled MassMigration
	err = json.Unmarshal(data, &unmarshalled)
	require.NoError(t, err)

	require.Equal(t, massMigration, unmarshalled)
}

func TestMassMigration_MarshalJSON(t *testing.T) {
	massMigration := MassMigration{
		FromStateID: ref.Uint32(1),
		SpokeID:     ref.Uint32(5),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &exampleSignature,
	}
	data, err := json.Marshal(massMigration)
	require.NoError(t, err)

	expected := fmt.Sprintf(
		`{"Type":"MASS_MIGRATION","FromStateID":1,"SpokeID":5,"Amount":"50","Fee":"10","Nonce":"0","Signature":%q}`,
		exampleSignatureHex,
	)
	require.Equal(t, expected, string(data))
}
