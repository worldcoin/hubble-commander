package models

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const PublicKeyLength = 128

var (
	ZeroPublicKey             = PublicKey{}
	ErrInvalidPublicKeyLength = errors.New("invalid public key length")
	publicKeyT                = reflect.TypeOf(PublicKey{})
)

type PublicKey [PublicKeyLength]byte

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

func (p *PublicKey) SetBytes(b []byte) error {
	if len(b) != PublicKeyLength {
		// TODO: wrap with stacktrace
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
