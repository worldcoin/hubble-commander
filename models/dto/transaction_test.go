package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	enumerr "github.com/Worldcoin/hubble-commander/models/enums/errors"
	"github.com/stretchr/testify/require"
)

func TestTransaction_UnmarshalJSON_Transfer(t *testing.T) {
	input := fmt.Sprintf(`{
			"Type":"TRANSFER",
			"FromStateID":1,
			"ToStateID":2,
			"Amount":"50",
			"Fee":"10",
			"Nonce":"0",
			"Signature":"%s"
		}`,
		"0x010203"+strings.Repeat("0", models.SignatureLength*2-6))
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.NoError(t, err)
	require.NotNil(t, tx.Parsed)
	require.IsType(t, Transfer{}, tx.Parsed)
	require.Equal(t, models.NewUint256(50), tx.Parsed.(Transfer).Amount)
	require.Equal(t, models.Signature{1, 2, 3}, *tx.Parsed.(Transfer).Signature)
}

func TestTransaction_UnmarshalJSON_Create2Transfer(t *testing.T) {
	input := fmt.Sprintf(`{
			"Type":"CREATE2TRANSFER",
			"FromStateID":1,
			"ToStateID":2,
			"ToPublicKey":"%s",
			"Amount":"50",
			"Fee":"10",
			"Nonce":"0",
			"Signature":"%s"
		}`,
		"0x010203"+strings.Repeat("0", models.PublicKeyLength*2-6),
		"0x040506"+strings.Repeat("0", models.SignatureLength*2-6),
	)
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.NoError(t, err)
	require.NotNil(t, tx.Parsed)
	require.IsType(t, Create2Transfer{}, tx.Parsed)
	require.Equal(t, models.NewUint256(50), tx.Parsed.(Create2Transfer).Amount)
	require.Equal(t, models.PublicKey{1, 2, 3}, *tx.Parsed.(Create2Transfer).ToPublicKey)
	require.Equal(t, models.Signature{4, 5, 6}, *tx.Parsed.(Create2Transfer).Signature)
}

func TestTransaction_UnmarshalJSON_UnknownType(t *testing.T) {
	input := `{"Type":"UNSUPPORTED_TYPE","FromStateID":1,"ToStateID":2,"Amount":"50","Fee":"10","Nonce":"0","Signature":"0xdeadbeef"}`
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.Equal(t, enumerr.NewUnsupportedError("transaction type"), err)
}
