package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type Uint256 struct {
	big.Int
}

func MakeUint256(value int64) Uint256 {
	return Uint256{*big.NewInt(value)}
}

func MakeUint256FromBig(value big.Int) Uint256 {
	return Uint256{value}
}

func NewUint256(value int64) *Uint256 {
	uint256 := MakeUint256(value)
	return &uint256
}

// Scan implements Scanner for database/sql.
func (a *Uint256) Scan(src interface{}) error {
	value, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("can't scan %T into Uint256", src)
	}

	a.SetString(string(value), 10)
	return nil
}

// Value implements valuer for database/sql.
func (a Uint256) Value() (driver.Value, error) {
	return a.Text(10), nil
}
