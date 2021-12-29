package main

import "C"

import (
	"encoding/hex"
	"math/big"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

var placeholderDomain = bls.Domain{0x00, 0x00, 0x00, 0x00}

func main() {}

func parseWallet(privateKey *C.char, domain *bls.Domain) (*bls.Wallet, error) {
	privateKeyDecoded, err := hex.DecodeString(C.GoString(privateKey))
	if err != nil {
		return nil, err
	}

	wallet, err := bls.NewWallet(privateKeyDecoded, *domain)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func parseDomain(domain *C.char) (*bls.Domain, error) {
	domainDecoded, err := hex.DecodeString(C.GoString(domain))
	if err != nil {
		return nil, err
	}
	domainBls, err := bls.DomainFromBytes(domainDecoded)
	if err != nil {
		return nil, err
	}

	return domainBls, nil
}

func parsePublicKey(pubKey *C.char) (*models.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(C.GoString(pubKey))
	if err != nil {
		return nil, err
	}

	var publicKey models.PublicKey
	copy(publicKey[:], publicKeyBytes)

	return &publicKey, nil
}

//export NewWalletPrivateKey
func NewWalletPrivateKey() *C.char {
	wallet, err := bls.NewRandomWallet(placeholderDomain)
	if err != nil {
		return nil
	}

	privateKey, _ := wallet.Bytes()
	return C.CString(hex.EncodeToString(privateKey))
}

//export GetWalletPublicKey
func GetWalletPublicKey(privateKey *C.char) *C.char {
	wallet, err := parseWallet(privateKey, &placeholderDomain)
	if err != nil {
		return nil
	}

	publicKey := wallet.PublicKey()
	return C.CString(hex.EncodeToString(publicKey[:]))
}

//export SignTransfer
func SignTransfer(from, to C.uint, amount, fee, nonce, privateKey, domain *C.char) *C.char {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return nil
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return nil
	}

	amountBigInt := new(big.Int)
	amountBigInt.SetString(C.GoString(amount), 10)

	feeBigInt := new(big.Int)
	feeBigInt.SetString(C.GoString(fee), 10)

	nonceBigInt := new(big.Int)
	nonceBigInt.SetString(C.GoString(nonce), 10)

	transfer, err := api.SignTransfer(wallet, dto.Transfer{
		FromStateID: ref.Uint32(uint32(from)),
		ToStateID:   ref.Uint32(uint32(to)),
		Amount:      models.NewUint256FromBig(*amountBigInt),
		Fee:         models.NewUint256FromBig(*feeBigInt),
		Nonce:       models.NewUint256FromBig(*nonceBigInt),
	})
	if err != nil {
		return nil
	}

	return C.CString(hex.EncodeToString(transfer.Signature.Bytes()))
}

//export SignCreate2Transfer
func SignCreate2Transfer(from C.uint, toPubKey, amount, fee, nonce, privateKey, domain *C.char) *C.char {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return nil
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return nil
	}

	amountBigInt := new(big.Int)
	amountBigInt.SetString(C.GoString(amount), 10)

	feeBigInt := new(big.Int)
	feeBigInt.SetString(C.GoString(fee), 10)

	nonceBigInt := new(big.Int)
	nonceBigInt.SetString(C.GoString(nonce), 10)

	toPublicKey, err := parsePublicKey(toPubKey)
	if err != nil {
		return nil
	}

	transfer, err := api.SignCreate2Transfer(wallet, dto.Create2Transfer{
		FromStateID: ref.Uint32(uint32(from)),
		ToPublicKey: toPublicKey,
		Amount:      models.NewUint256FromBig(*amountBigInt),
		Fee:         models.NewUint256FromBig(*feeBigInt),
		Nonce:       models.NewUint256FromBig(*nonceBigInt),
	})
	if err != nil {
		return nil
	}

	return C.CString(hex.EncodeToString(transfer.Signature.Bytes()))
}

//export SignMessage
func SignMessage(message, privateKey, domain *C.char) *C.char {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return nil
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return nil
	}

	signature, err := wallet.Sign([]byte(C.GoString(message)))
	if err != nil {
		return nil
	}

	return C.CString(hex.EncodeToString(signature.Bytes()))
}

//export VerifySignedMessage
func VerifySignedMessage(message, signature, pubKey, domain *C.char) C.int {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return C.int(-1)
	}

	signatureDecoded, err := hex.DecodeString(C.GoString(signature))
	if err != nil {
		return C.int(-1)
	}

	signatureObj, err := bls.NewSignatureFromBytes(signatureDecoded, *domainBls)
	if err != nil {
		return C.int(-1)
	}

	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return C.int(-1)
	}

	res, err := signatureObj.Verify([]byte(C.GoString(message)), publicKey)
	if err != nil {
		return C.int(-1)
	}

	if res {
		return C.int(1)
	}
	return C.int(0)
}
