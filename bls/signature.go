package bls

import (
	"math/big"

	"github.com/kilic/bn254/bls"
)

type Signature struct {
	sig      *bls.Signature
	verifier *bls.BLSVerifier
}

func NewSignature(signature *bls.Signature, domain Domain) *Signature {
	return &Signature{
		sig:      signature,
		verifier: bls.NewBLSVerifier(domain[:]),
	}
}

func (s *Signature) Domain() [32]byte {
	var domain [32]byte
	copy(domain[:], s.verifier.Domain)
	return domain
}

func (s *Signature) Verify(message []byte, publicKey *PublicKey) (bool, error) {
	return s.verifier.Verify(message, s.sig, publicKey.key)
}

func (s *Signature) ToBigInts() [2]*big.Int {
	bytes := s.sig.ToBytes()
	return [2]*big.Int{
		new(big.Int).SetBytes(bytes[:32]),
		new(big.Int).SetBytes(bytes[32:]),
	}
}
