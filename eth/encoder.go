package eth

import (
	"log"

	"github.com/Worldcoin/hubble-commander/contracts/frontend/transfer"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

// nolint
var tUint256 abi.Type

func init() {
	t, err := abi.NewType("uint256", "", nil)
	if err != nil {
		log.Fatal("Failed to create tUint256 constant")
	}
	tUint256 = t
}

func EncodeTransfer(tx transfer.OffchainTransfer) ([]uint8, error) {
	arguments := abi.Arguments{
		{Name: "txType", Type: tUint256},
		{Name: "fromIndex", Type: tUint256},
		{Name: "toIndex", Type: tUint256},
		{Name: "amount", Type: tUint256},
		{Name: "fee", Type: tUint256},
		{Name: "nonce", Type: tUint256},
	}
	encodedBytes, err := arguments.Pack(
		tx.TxType,
		tx.FromIndex,
		tx.ToIndex,
		tx.Amount,
		tx.Fee,
		tx.Nonce,
	)
	if err != nil {
		return nil, err
	}
	return encodedBytes, nil
}
