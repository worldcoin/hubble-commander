package utils

import (
	"encoding/hex"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

func RandomBytes(size uint64) []byte {
	bytes := make([]byte, size)
	rand.Read(bytes)
	return bytes
}

func RandomHex(length uint64) string {
	return hex.EncodeToString(RandomBytes(length / 2))
}

func RandomHash() common.Hash {
	return common.BytesToHash(RandomBytes(32))
}

func RandomAddress() common.Address {
	return common.BytesToAddress(RandomBytes(20))
}

func RandomBigInt() *big.Int {
	return new(big.Int).SetUint64(rand.Uint64())
}
