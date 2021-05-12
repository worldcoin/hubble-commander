package bls

import (
	"crypto/rand"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/kilic/bn254/bls"
)

type (
	KeyPair = bls.KeyPair
)

type Wallet struct {
	signer bls.BLSSigner
}

func NewWallet(privateKey []byte, domain Domain) (*Wallet, error) {
	keyPair, err := bls.NewKeyPairFromSecret(privateKey)
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
	domain, _ := DomainFromBytes(w.signer.Domain)
	return *domain
}

func (w *Wallet) PublicKey() *models.PublicKey {
	return fromBLSPublicKey(w.signer.Account.Public)
}

func (w *Wallet) Bytes() (privateKey, publicKey []byte) {
	accountBytes := w.signer.Account.ToBytes()
	privateKey = accountBytes[128:]
	publicKey = fromBLSPublicKey(w.signer.Account.Public).Bytes()
	return
}
