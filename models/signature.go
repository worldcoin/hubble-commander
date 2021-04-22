package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils"
)

type Signature [64]byte

func MakeSignature(first, second int64) Signature {
	var signature [64]byte
	copy(signature[:32], utils.PadLeft(big.NewInt(first).Bytes(), 32))
	copy(signature[32:], utils.PadLeft(big.NewInt(second).Bytes(), 32))
	return signature
}

func (s Signature) Bytes() []byte {
	return s[:]
}

func (s *Signature) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Signature", src)
	}
	if len(value) != 64 {
		return fmt.Errorf("invalid signature length")
	}

	copy(s[:], value)
	return nil
}

func (s Signature) Value() (driver.Value, error) {
	return s.Bytes(), nil
}

func (s *Signature) ToBigIntPointers() [2]*big.Int {
	return [2]*big.Int{
		new(big.Int).SetBytes(s[:32]),
		new(big.Int).SetBytes(s[32:]),
	}
}
