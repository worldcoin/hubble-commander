package models

import (
	"database/sql/driver"
	"fmt"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const PublicKeyLength = 128

type PublicKey [PublicKeyLength]byte

var publicKeyT = reflect.TypeOf(PublicKey{})

func MakePublicKeyFromInts(ints [4]*big.Int) PublicKey {
	publicKey := PublicKey{}
	copy(publicKey[:32], utils.PadLeft(ints[0].Bytes(), 32))
	copy(publicKey[32:64], utils.PadLeft(ints[1].Bytes(), 32))
	copy(publicKey[64:96], utils.PadLeft(ints[2].Bytes(), 32))
	copy(publicKey[96:], utils.PadLeft(ints[3].Bytes(), 32))
	return publicKey
}

// nolint:gocritic
func (p PublicKey) Bytes() []byte {
	return p[:]
}

func (p *PublicKey) BigInts() [4]*big.Int {
	return [4]*big.Int{
		new(big.Int).SetBytes(p[:32]),
		new(big.Int).SetBytes(p[32:64]),
		new(big.Int).SetBytes(p[64:96]),
		new(big.Int).SetBytes(p[96:]),
	}
}

func (p *PublicKey) String() string {
	return hexutil.Encode(p[:])
}

func (p *PublicKey) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(srcB) != PublicKeyLength {
		return fmt.Errorf("can't scan []byte of len %d into PublicKey, want %d", len(srcB), PublicKeyLength)
	}
	copy(p[:], srcB)
	return nil
}

// nolint:gocritic
func (p PublicKey) Value() (driver.Value, error) {
	return p[:], nil
}

func (p *PublicKey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(publicKeyT, input, p[:])
}

// nolint:gocritic
func (p PublicKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(p[:]).MarshalText()
}
