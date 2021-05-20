package models

import (
	"database/sql/driver"
	"encoding/json"
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

func NewUint256FromBig(value big.Int) *Uint256 {
	uint256 := MakeUint256FromBig(value)
	return &uint256
}

func (u *Uint256) Add(other *Uint256) *Uint256 {
	sum := new(big.Int).Add(&u.Int, &other.Int)
	return &Uint256{*sum}
}

func (u *Uint256) Sub(other *Uint256) *Uint256 {
	diff := new(big.Int).Sub(&u.Int, &other.Int)
	return &Uint256{*diff}
}

func (u *Uint256) Mul(other *Uint256) *Uint256 {
	product := new(big.Int).Mul(&u.Int, &other.Int)
	return &Uint256{*product}
}

func (u *Uint256) Div(other *Uint256) *Uint256 {
	quotient := new(big.Int).Div(&u.Int, &other.Int)
	return &Uint256{*quotient}
}

func (u *Uint256) AddN(other int64) *Uint256 {
	return u.Add(NewUint256(other))
}

func (u *Uint256) SubN(other int64) *Uint256 {
	return u.Sub(NewUint256(other))
}

func (u *Uint256) MulN(other int64) *Uint256 {
	return u.Mul(NewUint256(other))
}

func (u *Uint256) DivN(other int64) *Uint256 {
	return u.Div(NewUint256(other))
}

func (u *Uint256) Cmp(other *Uint256) int {
	return u.Int.Cmp(&other.Int)
}

func (u *Uint256) CmpN(other int64) int {
	return u.Cmp(NewUint256(other))
}

// Scan implements Scanner for database/sql.
func (u *Uint256) Scan(src interface{}) error {
	value, ok := src.([]uint8)
	if !ok {
		return fmt.Errorf("can't scan %T into Uint256", src)
	}

	u.SetString(string(value), 10)
	return nil
}

// Value implements valuer for database/sql.
func (u Uint256) Value() (driver.Value, error) {
	return u.Text(10), nil
}

func (u Uint256) MarshalJSON() ([]byte, error) {
	jsonText, err := json.Marshal(u.Text(10))
	if err != nil {
		return nil, err
	}
	return jsonText, nil
}

func (u *Uint256) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	_, ok := u.SetString(str, 10)

	if !ok {
		return fmt.Errorf("error unmarshaling Uint256")
	}

	return nil
}

func (u *Uint256) SetBytes(data []byte) {
	u.Int.SetBytes(data)
	if u.CmpN(0) == 0 {
		*u = MakeUint256(0)
	}
}
