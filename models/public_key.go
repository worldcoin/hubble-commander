package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/utils"
)

type PublicKey [128]byte

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
	return hex.EncodeToString(p.Bytes())
}

func (p *PublicKey) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(value) != 128 {
		return fmt.Errorf("invalid public key length")
	}

	copy(p[:], value)
	return nil
}

// nolint:gocritic
func (p PublicKey) Value() (driver.Value, error) {
	return p.Bytes(), nil
}

func (p *PublicKey) UnmarshalJSON(b []byte) error {
	decodedBytes, err := hex.Decode(p[:], b[1:len(b)-1])
	if err != nil {
		return err
	}
	if decodedBytes != 128 {
		return fmt.Errorf("invalid public key")
	}

	return nil
}

// nolint:gocritic
func (p PublicKey) MarshalJSON() ([]byte, error) {
	marshalizedPublicKey, err := json.Marshal(hex.EncodeToString(p[:]))
	if err != nil {
		return nil, err
	}
	return marshalizedPublicKey, nil
}
