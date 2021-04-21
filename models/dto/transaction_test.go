package dto

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/stretchr/testify/require"
)

func TestTransaction_UnmarshalJSON_Transfer(t *testing.T) {
	input := `{
		"Type":"TRANSFER",
		"FromStateID":1,
		"ToStateID":2,
		"Amount":"50",
		"Fee":"10",
		"Nonce":"0",
		"Signature":"0xdeadbeef"
	}`

	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.NoError(t, err)
	require.NotNil(t, tx.Parsed)
	require.IsType(t, Transfer{}, tx.Parsed)
	require.Equal(t, models.NewUint256(50), tx.Parsed.(Transfer).Amount)
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
			"Signature":"0xdeadbeef"
		}`,
		"010203"+strings.Repeat("0", 250),
	)
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.NoError(t, err)
	require.NotNil(t, tx.Parsed)
	require.IsType(t, Create2Transfer{}, tx.Parsed)
	require.Equal(t, models.NewUint256(50), tx.Parsed.(Create2Transfer).Amount)
	require.Equal(t, models.PublicKey{1, 2, 3}, *tx.Parsed.(Create2Transfer).ToPublicKey)
}

func TestTransaction_UnmarshalJSON_UnknownType(t *testing.T) {
	input := `{"Type":"UNSUPPORTED_TYPE","FromStateID":1,"ToStateID":2,"Amount":"50","Fee":"10","Nonce":"0","Signature":"0xdeadbeef"}`
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.Equal(t, txtype.ErrUnsupportedTransactionType, err)
}
