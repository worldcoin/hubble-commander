package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils"
)

type Signature [64]byte

func MakeRandomSignature() Signature {
	var signature Signature
	copy(signature[:], utils.RandomBytes(64))
	return signature
}

func (s Signature) Bytes() []byte {
	return s[:]
}

func (s *Signature) BigInts() [2]*big.Int {
	return [2]*big.Int{
		new(big.Int).SetBytes(s[:32]),
		new(big.Int).SetBytes(s[32:]),
	}
}

func (s *Signature) String() string {
	return hex.EncodeToString(s.Bytes())
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

func (s *Signature) UnmarshalJSON(b []byte) error {
	decodedBytes, err := hex.Decode(s[:], b[1:len(b)-1])
	if err != nil {
		return err
	}
	if decodedBytes != 64 {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(s[:]))
}
