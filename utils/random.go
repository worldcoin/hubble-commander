package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
)

func SafeRandomBytes(size uint64) ([]byte, error) {
	bytes := make([]byte, size)
	_, err := crand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func RandomBytes(size uint64) []byte {
	bytes := make([]byte, size)
	// nolint:gosec
	rand.Read(bytes)
	return bytes
}

func RandomHex(length uint64) string {
	return hex.EncodeToString(RandomBytes(length / 2))
}

func RandomHash() common.Hash {
	return common.BytesToHash(RandomBytes(32))
}

func NewRandomHash() *common.Hash {
	newHash := RandomHash()
	return &newHash
}

func RandomAddress() common.Address {
	return common.BytesToAddress(RandomBytes(20))
}

func RandomBigInt() *big.Int {
	// nolint:gosec
	return new(big.Int).SetUint64(rand.Uint64())
}
