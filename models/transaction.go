package models

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Hash      common.Hash
	FromIndex *big.Int
	ToIndex   *big.Int
	Amount    *big.Int
	Fee       *big.Int
	Nonce     *big.Int
	// TODO: Right now decoder expects a base64 string here, we could define a custom type with interface implementation to expect a hex string
	Signature []byte
}
