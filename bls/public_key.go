package bls

import (
	"math/big"

	"github.com/kilic/bn254/bls"
)

type PublicKey struct {
	key *bls.PublicKey
}

func (p *PublicKey) ToBigInts() [4]*big.Int {
	bytes := p.key.ToBytes()
	return [4]*big.Int{
		new(big.Int).SetBytes(bytes[:32]),
		new(big.Int).SetBytes(bytes[32:64]),
		new(big.Int).SetBytes(bytes[64:96]),
		new(big.Int).SetBytes(bytes[96:]),
	}
}
