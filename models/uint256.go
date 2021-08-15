package models

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var ErrUnmarshalUint256 = errors.New("error unmarshalling Uint256")

type Uint256 struct {
	uint256.Int
}

func MakeUint256(value uint64) Uint256 {
	res := Uint256{}
	res.Int.SetUint64(value)
	return res
}

func MakeUint256FromBig(value big.Int) Uint256 {
	res := Uint256{}
	overflow := res.Int.SetFromBig(&value)
	if overflow {
		panic("overflow occurred")
	}
	return res
}

func NewUint256(value uint64) *Uint256 {
	newUint256 := MakeUint256(value)
	return &newUint256
}

func NewUint256FromBig(value big.Int) *Uint256 {
	newUint256 := MakeUint256FromBig(value)
	return &newUint256
}

func (u *Uint256) Add(other *Uint256) *Uint256 {
	sum := uint256.NewInt(0).Add(&u.Int, &other.Int)
	return &Uint256{*sum}
}

func (u *Uint256) Sub(other *Uint256) *Uint256 {
	diff := uint256.NewInt(0).Sub(&u.Int, &other.Int)
	return &Uint256{*diff}
}

func (u *Uint256) Mul(other *Uint256) *Uint256 {
	product := uint256.NewInt(0).Mul(&u.Int, &other.Int)
	return &Uint256{*product}
}

func (u *Uint256) Div(other *Uint256) *Uint256 {
	quotient := uint256.NewInt(0).Div(&u.Int, &other.Int)
	return &Uint256{*quotient}
}

func (u *Uint256) AddN(other uint64) *Uint256 {
	return u.Add(NewUint256(other))
}

func (u *Uint256) SubN(other uint64) *Uint256 {
	return u.Sub(NewUint256(other))
}

func (u *Uint256) MulN(other uint64) *Uint256 {
	return u.Mul(NewUint256(other))
}

func (u *Uint256) DivN(other uint64) *Uint256 {
	return u.Div(NewUint256(other))
}

func (u *Uint256) Cmp(other *Uint256) int {
	return u.Int.Cmp(&other.Int)
}

func (u *Uint256) CmpN(other uint64) int {
	return u.Cmp(NewUint256(other))
}

// Scan implements Scanner for database/sql.
func (u *Uint256) Scan(src interface{}) error {
	errorMessage := "can't scan %T into Uint256"

	value, ok := src.([]uint8)
	if !ok {
		return errors.Errorf(errorMessage, src)
	}
	bigValue, ok := new(big.Int).SetString(string(value), 10)
	if !ok {
		return errors.Errorf(errorMessage, src)
	}

	overflow := u.Int.SetFromBig(bigValue)
	if overflow {
		return errors.Errorf(errorMessage, src)
	}

	return nil
}

// Value implements valuer for database/sql.
func (u Uint256) Value() (driver.Value, error) {
	return u.Int.ToBig().Text(10), nil
}

func (u *Uint256) String() string {
	return u.Int.ToBig().Text(10)
}

func (u Uint256) MarshalJSON() ([]byte, error) {
	jsonText, err := json.Marshal(u.Int.ToBig().Text(10))
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

	bigValue, ok := new(big.Int).SetString(str, 10)
	if !ok {
		return ErrUnmarshalUint256
	}

	overflow := u.Int.SetFromBig(bigValue)
	if overflow {
		return ErrUnmarshalUint256
	}

	return nil
}

func (u Uint256) MarshalYAML() (interface{}, error) {
	return "0x" + u.Int.ToBig().Text(16), nil
}

func (u *Uint256) UnmarshalYAML(b []byte) error {
	var number string
	err := yaml.Unmarshal(b, &number)
	if err != nil {
		return ErrUnmarshalUint256
	}

	return nil
}

func (u *Uint256) SetBytes(data []byte) {
	u.Int.SetBytes(data)
}
