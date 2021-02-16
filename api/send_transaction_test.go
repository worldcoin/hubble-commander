package api

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestApi_SendTransaction(t *testing.T) {
	api := Api{&config.Config{Version: "v0123"}}
	tx := IncomingTransaction{
		FromIndex: big.NewInt(1),
		ToIndex:   big.NewInt(2),
		Amount:    big.NewInt(50),
		Fee:       big.NewInt(10),
		Nonce:     big.NewInt(0),
		Signature: []byte{97, 100, 115, 97, 100, 115, 97, 115, 100, 97, 115, 100},
	}
	hash, err := api.SendTransaction(tx)
	require.NoError(t, err)
	require.Equal(t, common.HexToHash("0x3e136a19201d6fc73c4e3c76951edfb94eb9a7a0c7e15492696ffddb3e1b2c68"), hash)
}
