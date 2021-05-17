package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const SignatureLength = 64

type Signature [SignatureLength]byte

var signatureT = reflect.TypeOf(Signature{})

func MakeRandomSignature() Signature {
	var signature Signature
	copy(signature[:], utils.RandomBytes(64))
	return signature
}

func NewRandomSignature() *Signature {
	signature := MakeRandomSignature()
	return &signature
}

func MakeSignatureFromBigInts(ints []*big.Int) Signature {
	var signature Signature

	copy(signature[0:32], ints[0].Bytes())
	copy(signature[32:64], ints[1].Bytes())

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
	return hexutil.Encode(s[:])
}

func (s *Signature) Scan(src interface{}) error {
	srcBytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Signature", src)
	}
	if len(srcBytes) != SignatureLength {
		return fmt.Errorf("can't scan []byte of len %d into Signature, want %d", len(srcBytes), SignatureLength)
	}
	copy(s[:], srcBytes)
	return nil
}

func (s Signature) Value() (driver.Value, error) {
	return s[:], nil
}

func (s *Signature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(signatureT, input, s[:])
}

func (s Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(s[:]).MarshalText()
}
