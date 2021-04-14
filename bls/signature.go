package bls

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
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

func NewSignatureFromBytes(signatureBytes []byte, domain Domain) (*Signature, error) {
	signature, err := bls.SignatureFromBytes(signatureBytes)
	if err != nil {
		return nil, err
	}
	return NewSignature(signature, domain), nil
}

func (s *Signature) Domain() [32]byte {
	var domain [32]byte
	copy(domain[:], s.verifier.Domain)
	return domain
}

func (s *Signature) Verify(message []byte, publicKey *models.PublicKey) (bool, error) {
	return s.verifier.Verify(message, s.sig, toBLSPublicKey(publicKey))
}

func (s *Signature) BigInts() [2]*big.Int {
	bytes := s.sig.ToBytes()
	return [2]*big.Int{
		new(big.Int).SetBytes(bytes[:32]),
		new(big.Int).SetBytes(bytes[32:]),
	}
}

func (s *Signature) Bytes() []byte {
	return s.sig.ToBytes()
}
