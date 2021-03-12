package models

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
)

type PublicKey [128]byte

// nolint:gocritic
func (p PublicKey) Bytes() []byte {
	return p[:]
}

func (p *PublicKey) IntArray() [4]*big.Int {
	ints := [4]*big.Int{}

	ints[0] = new(big.Int).SetBytes(p[0:32])
	ints[1] = new(big.Int).SetBytes(p[32:64])
	ints[2] = new(big.Int).SetBytes(p[64:96])
	ints[3] = new(big.Int).SetBytes(p[96:128])

	return ints
}

func (p *PublicKey) String() string {
	return hex.EncodeToString(p.Bytes())
}

func MakePublicKeyFromUint256(ints [4]*big.Int) PublicKey {
	publicKey := PublicKey{}
	copy(publicKey[0:32], ints[0].Bytes())
	copy(publicKey[32:64], ints[1].Bytes())
	copy(publicKey[64:86], ints[2].Bytes())
	copy(publicKey[96:128], ints[3].Bytes())

	return publicKey
}

func (p *PublicKey) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into PublicKey", src)
	}
	if len(value) != 128 {
		return fmt.Errorf("invalid signature length")
	}

	copy(p[:], value)
	return nil
}

// nolint:gocritic
// Value implements valuer for database/sql.
func (p PublicKey) Value() (driver.Value, error) {
	return p.Bytes(), nil
}
