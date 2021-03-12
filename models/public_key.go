package models

import (
	"database/sql/driver"
	"fmt"
)

type PublicKey [128]byte

func (p PublicKey) Bytes() []byte {
	return p[:]
}

func (p *PublicKey) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(value) != 128 {
		return fmt.Errorf("invalid signature length")
	}

	copy(p[:], value)
	return nil
}

// Value implements valuer for database/sql.
func (p PublicKey) Value() (driver.Value, error) {
	return p.Bytes(), nil
}
