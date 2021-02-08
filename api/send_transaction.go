package api

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

type IncomingTransaction struct {
	FromIndex *big.Int
	ToIndex   *big.Int
	Amount    *big.Int
	Fee       *big.Int
	Nonce     *big.Int
	// TODO: Right now decoder expects a base64 string here, we could define a custom type with interface implementation to expect a hex string
	Signature []byte
}

func (a *Api) SendTransaction(tx IncomingTransaction) (h common.Hash, err error) {
	fmt.Printf("%+v\n", tx)
	h, err = rlpHash(tx)
	if err != nil {
		return
	}
	return h, nil
}

func rlpHash(x interface{}) (h common.Hash, err error) {
	hw := sha3.NewLegacyKeccak256()
	if err = rlp.Encode(hw, x); err != nil {
		return
	}
	hw.Sum(h[:0])
	return h, nil
}
