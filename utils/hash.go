package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func HashTwo(a, b common.Hash) common.Hash {
	buf := make([]byte, 64)
	copy(buf[0:32], a.Bytes())
	copy(buf[32:64], b.Bytes())
	return crypto.Keccak256Hash(buf)
}

/*
 * This syntax was added in go 1.17
 * - https://tip.golang.org/ref/spec#Conversions_from_slice_to_array_pointer
 */
func HashToByteArray(a *common.Hash) [32]byte {
	return *(*[32]byte)(a.Bytes())
}
