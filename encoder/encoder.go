package encoder

import (
	"github.com/Worldcoin/hubble-commander/contracts/frontend/generic"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var tUint256, _ = abi.NewType("uint256", "", nil)

func EncodeUserState(state generic.TypesUserState) ([]byte, error) {
	arguments := abi.Arguments{
		{Name: "pubkeyID", Type: tUint256},
		{Name: "tokenID", Type: tUint256},
		{Name: "balance", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	return arguments.Pack(
		state.PubkeyID,
		state.TokenID,
		state.Balance,
		state.Nonce,
	)
}
