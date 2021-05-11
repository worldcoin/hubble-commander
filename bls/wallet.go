package bls

import (
	"crypto/rand"
	"errors"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/kilic/bn254/bls"
)

type (
	KeyPair = bls.KeyPair
	Domain  = [32]byte
)

var (
	testDomain = Domain{0x00, 0x00, 0x00, 0x00}

	ErrInvalidDomainLength = errors.New("invalid domain length")
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

func DomainFromBytes(data []byte) (*Domain, error) {
	if len(data) != 32 {
		return nil, ErrInvalidDomainLength
	}
	var domain [32]byte
	copy(domain[:], data)
	return &domain, nil
}
