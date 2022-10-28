package models

import (
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

func MakeSignatureFromBigInts(ints [2]*big.Int) Signature {
	var signature Signature

	copy(signature[0:32], utils.PadLeft(ints[0].Bytes(), 32))
	copy(signature[32:64], utils.PadLeft(ints[1].Bytes(), 32))

	return signature
}

func (s Signature) Bytes() []byte {
	return s[:]
}

func (s *Signature) SetBytes(data []byte) error {
	if len(data) != SignatureLength {
		// TODO: errors.WithStack
		//       can't do it now because of the possibility that a caller is failing to use
		//       errors.Is instead of strict equality
		return ErrInvalidLength
	}

	copy(s[:], data)
	return nil
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

func (s *Signature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(signatureT, input, s[:])
}

func (s Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(s[:]).MarshalText()
}
