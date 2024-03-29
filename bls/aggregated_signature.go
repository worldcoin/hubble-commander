package bls

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/kilic/bn254/bls"
)

type AggregatedSignature struct {
	*Signature
}

func NewAggregatedSignature(signatures []*Signature) *AggregatedSignature {
	if len(signatures) == 0 {
		panic("signatures slice cannot be empty")
	}
	domain := signatures[0].Domain()
	blsSignatures := signaturesToBls(signatures, domain)
	return NewAggregatedSignatureFromBls(blsSignatures, domain)
}

func NewAggregatedSignatureFromBls(signatures []*bls.Signature, domain Domain) *AggregatedSignature {
	return &AggregatedSignature{
		Signature: NewSignature(bls.AggregateSignatures(signatures), domain),
	}
}

func (s *AggregatedSignature) Verify(messages [][]byte, publicKeys []*models.PublicKey) (bool, error) {
	keys := make([]*bls.PublicKey, 0, len(publicKeys))
	for _, pk := range publicKeys {
		keys = append(keys, toBLSPublicKey(pk))
	}
	return s.verifier.VerifyAggregate(messages, keys, s.sig)
}

func signaturesToBls(signatures []*Signature, domain Domain) []*bls.Signature {
	blsSignatures := make([]*bls.Signature, 0, len(signatures))
	for _, s := range signatures {
		if s.Domain() != domain {
			panic("all signatures must have the same domain")
		}
		blsSignatures = append(blsSignatures, s.sig)
	}
	return blsSignatures
}
