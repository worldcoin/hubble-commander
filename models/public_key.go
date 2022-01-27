package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const PublicKeyLength = 128

var ErrInvalidPublicKeyLength = errors.New("invalid public key length")

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

func NewRandomPublicKey() (*PublicKey, error) {
	var publicKey PublicKey
	bytes, err := utils.SafeRandomBytes(PublicKeyLength)
	if err != nil {
		return nil, err
	}
	copy(publicKey[:], bytes)
	return &publicKey, nil
}

// nolint:gocritic
func (p PublicKey) Bytes() []byte {
	return p[:]
}

func (p *PublicKey) SetBytes(b []byte) error {
	if len(b) != PublicKeyLength {
		return ErrInvalidPublicKeyLength
	}

	copy(p[0:], b)
	return nil
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
	srcBytes, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(srcBytes) != PublicKeyLength {
		return fmt.Errorf("can't scan []byte of len %d into PublicKey, want %d", len(srcBytes), PublicKeyLength)
	}
	copy(p[:], srcBytes)
	return nil
}

// nolint:gocritic
func (p PublicKey) Value() (driver.Value, error) {
	return p[:], nil
}

func (p *PublicKey) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(publicKeyT, input, p[:])
}

func (p *PublicKey) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var publicKey string
	err := unmarshal(&publicKey)
	if err != nil {
		return err
	}
	decodedHex, err := hexutil.Decode(publicKey)
	if err != nil {
		return err
	}
	if len(decodedHex) != PublicKeyLength {
		return ErrInvalidPublicKeyLength
	}
	copy(p[:], decodedHex)
	return err
}

// nolint:gocritic
func (p PublicKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(p[:]).MarshalText()
}
