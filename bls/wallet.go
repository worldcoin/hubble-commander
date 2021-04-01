package bls

import (
	"crypto/rand"

	"github.com/kilic/bn254/bls"
)

type (
	KeyPair = bls.KeyPair
	Domain  = [32]byte
)

var testDomain = Domain{0x00, 0x00, 0x00, 0x00}

type Wallet struct {
	signer bls.BLSSigner
}

func NewWallet(secretKey []byte, domain Domain) (*Wallet, error) {
	keyPair, err := bls.NewKeyPairFromSecret(secretKey)
	if err != nil {
		return nil, err
	}
	return NewWalletFromKeyPair(keyPair, domain), nil
}

func NewRandomWallet(domain Domain) (*Wallet, error) {
	keyPair, err := bls.NewKeyPair(rand.Reader)
	if err != nil {
		return nil, err
	}
	return NewWalletFromKeyPair(keyPair, domain), nil
}

func NewWalletFromKeyPair(account *KeyPair, domain Domain) *Wallet {
	signer := bls.BLSSigner{
		Account: account,
		Domain:  domain[:],
	}
	return &Wallet{signer: signer}
}

func (w *Wallet) Sign(data []byte) (*Signature, error) {
	signature, err := w.signer.Sign(data)
	if err != nil {
		return nil, err
	}
	return NewSignature(signature, w.Domain()), nil
}

func (w *Wallet) Domain() Domain {
	var domain [32]byte
	copy(domain[:], w.signer.Domain)
	return domain
}

func (w *Wallet) PublicKey() *PublicKey {
	return &PublicKey{key: w.signer.Account.Public}
}

func (w *Wallet) Bytes() (secretKey, publicKey []byte) {
	accountBytes := w.signer.Account.ToBytes()
	return accountBytes[128:], accountBytes[:128]
}
