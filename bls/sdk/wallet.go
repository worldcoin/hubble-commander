package HubbleSDK

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

var placeholderDomain = bls.Domain{0x00, 0x00, 0x00, 0x00}

func parseWallet(privateKey string, domain *bls.Domain) (*bls.Wallet, error) {
	privateKeyDecoded, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	wallet, err := bls.NewWallet(privateKeyDecoded, *domain)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func parseDomain(domain string) (*bls.Domain, error) {
	domainDecoded, err := hex.DecodeString(domain)
	if err != nil {
		return nil, err
	}
	domainBls, err := bls.DomainFromBytes(domainDecoded)
	if err != nil {
		return nil, err
	}

	return domainBls, nil
}

func parsePublicKey(pubKey string) (*models.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	var publicKey models.PublicKey
	copy(publicKey[:], publicKeyBytes)

	return &publicKey, nil
}

func parseUint256(uint256 string) (*models.Uint256, error) {
	value := new(big.Int)
	if _, success := value.SetString(uint256, 10); !success {
		return nil, errors.New("failed to parse Uint256")
	}
	return models.NewUint256FromBig(*value), nil
}

//export NewWalletPrivateKey
func NewWalletPrivateKey() (string, error) {
	wallet, err := bls.NewRandomWallet(placeholderDomain)
	if err != nil {
		return "", err
	}

	privateKey, _ := wallet.Bytes()
	return hex.EncodeToString(privateKey), nil
}

//export GetWalletPublicKey
func GetWalletPublicKey(privateKey string) (string, error) {
	wallet, err := parseWallet(privateKey, &placeholderDomain)
	if err != nil {
		return "", err
	}

	publicKey := wallet.PublicKey()
	return hex.EncodeToString(publicKey[:]), nil
}

//export SignTransfer
func SignTransfer(from, to int, amount, fee, nonce, privateKey, domain string) (string, error) {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return "", err
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return "", err
	}

	amountUint256, err := parseUint256(amount)
	if err != nil {
		return "", err
	}
	feeUint256, err := parseUint256(fee)
	if err != nil {
		return "", err
	}
	nonceUint256, err := parseUint256(nonce)
	if err != nil {
		return "", err
	}

	transfer, err := api.SignTransfer(wallet, dto.Transfer{
		FromStateID: ref.Uint32(uint32(from)),
		ToStateID:   ref.Uint32(uint32(to)),
		Amount:      amountUint256,
		Fee:         feeUint256,
		Nonce:       nonceUint256,
	})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(transfer.Signature.Bytes()), nil
}

//export SignCreate2Transfer
func SignCreate2Transfer(from int, toPubKey, amount, fee, nonce, privateKey, domain string) (string, error) {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return "", err
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return "", err
	}

	amountUint256, err := parseUint256(amount)
	if err != nil {
		return "", err
	}
	feeUint256, err := parseUint256(fee)
	if err != nil {
		return "", err
	}
	nonceUint256, err := parseUint256(nonce)
	if err != nil {
		return "", err
	}

	toPublicKey, err := parsePublicKey(toPubKey)
	if err != nil {
		return "", err
	}

	transfer, err := api.SignCreate2Transfer(wallet, dto.Create2Transfer{
		FromStateID: ref.Uint32(uint32(from)),
		ToPublicKey: toPublicKey,
		Amount:      amountUint256,
		Fee:         feeUint256,
		Nonce:       nonceUint256,
	})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(transfer.Signature.Bytes()), nil
}

//export SignMessage
func SignMessage(message, privateKey, domain string) (string, error) {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return "", err
	}

	wallet, err := parseWallet(privateKey, domainBls)
	if err != nil {
		return "", err
	}

	signature, err := wallet.Sign([]byte(message))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature.Bytes()), nil
}

//export VerifySignedMessage
func VerifySignedMessage(message, signature, pubKey, domain string) int {
	domainBls, err := parseDomain(domain)
	if err != nil {
		return -1
	}

	signatureDecoded, err := hex.DecodeString(signature)
	if err != nil {
		return -1
	}

	signatureObj, err := bls.NewSignatureFromBytes(signatureDecoded, *domainBls)
	if err != nil {
		return -1
	}

	publicKey, err := parsePublicKey(pubKey)
	if err != nil {
		return -1
	}

	res, err := signatureObj.Verify([]byte(message), publicKey)
	if err != nil {
		return -1
	}

	if res {
		return 1
	}
	return 0
}
