package bls

import (
	"crypto/rand"

	"github.com/kilic/bn254/bls"
)

var DefaultDomain = [32]byte{0x00, 0x00, 0x00, 0x00}

type Wallet struct {
	signer bls.BLSSigner
}

func BytesToSignature(b []byte) (*bls.Signature, error) {
	sig, err := bls.SignatureFromBytes(b)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func NewWallet(domain [32]byte) (*Wallet, error) {
	newAccount, err := bls.NewKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}
	signer := bls.BLSSigner{Account: newAccount, Domain: domain[:]}
	return &Wallet{signer: signer}, nil
}

func (w *Wallet) Bytes() (secretKey, pubkey []byte) {
	accountBytes := w.signer.Account.ToBytes()
	secretBytes := accountBytes[128:]
	pubkeyBytes := accountBytes[:128]
	return secretBytes, pubkeyBytes
}

func SecretToWallet(secretKey []byte, domain [32]byte) (*Wallet, error) {
	keyPair, err := bls.NewKeyPairFromSecret(secretKey)
	if err != nil {
		return nil, err
	}
	signer := bls.BLSSigner{Account: keyPair, Domain: domain[:]}
	return &Wallet{signer: signer}, nil
}

func (w *Wallet) Sign(data []byte) (*bls.Signature, error) {
	signature, err := w.signer.Sign(data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (w *Wallet) VerifySignature(signature *bls.Signature, data []byte, pubkey *bls.PublicKey) (bool, error) {
	verifier := bls.NewBLSVerifier(w.signer.Domain)
	valid, err := verifier.Verify(data, signature, pubkey)
	return valid, err
}

func VerifyAggregatedSignature(
	aggregateSignature bls.Signature,
	data []bls.Message,
	pubkeys []*bls.PublicKey,
	domain [32]byte,
) (bool, error) {
	verifier := bls.NewBLSVerifier(domain[:])
	return verifier.VerifyAggregate(data, pubkeys, &aggregateSignature)
}

func NewAggregateSignature(signatures []*bls.Signature) bls.Signature {
	aggregatedSig := bls.AggregateSignatures(signatures)
	return *aggregatedSig
}
