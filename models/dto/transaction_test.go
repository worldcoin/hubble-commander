package dto

import (
	"encoding/json"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func TestTransaction_UnmarshalJSON_Transfer(t *testing.T) {
	input := `{"Type":1,"FromStateID":1,"ToStateID":2,"Amount":"50","Fee":"10","Nonce":"0","Signature":"0xdeadbeef"}` // nolint:goconst
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.NoError(t, err)
	require.NotNil(t, tx.Parsed)
	require.IsType(t, Transfer{}, tx.Parsed)
	require.Equal(t, tx.Parsed.(Transfer).Amount, models.NewUint256(50))
}

func TestTransaction_UnmarshalJSON_UnknownType(t *testing.T) {
	input := `{"Type":123,"FromStateID":1,"ToStateID":2,"Amount":"50","Fee":"10","Nonce":"0","Signature":"0xdeadbeef"}`
	var tx Transaction
	err := json.Unmarshal([]byte(input), &tx)
	require.Equal(t, ErrUnsupportedType, err)
}
