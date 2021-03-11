package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

// TODO: Consider representing this as a 64 byte array instead
type Signature [2]Uint256

func (s *Signature) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Signature", src)
	}
	if len(value) != 64 {
		return fmt.Errorf("invalid signature length")
	}

	s[0].SetBytes(value[0:32])
	s[1].SetBytes(value[32:64])
	return nil
}

// Value implements valuer for database/sql.
func (s Signature) Value() (driver.Value, error) {
	buf := make([]byte, 0, 64)

	buf = append(buf, s[0].Bytes()...)
	buf = append(buf, s[1].Bytes()...)

	return buf, nil
}

func (s *Signature) ToBigIntPointers() [2]*big.Int {
	return [2]*big.Int{&s[0].Int, &s[1].Int}
}
