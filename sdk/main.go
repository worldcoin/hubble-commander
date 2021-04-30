package main

import "C"

import (
	"encoding/hex"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func main() {}

func parseWallet(privateKey *C.char) (*bls.Wallet, error) {
	privateKeyDecoded, err := hex.DecodeString(C.GoString(privateKey))
	if err != nil {
		return nil, err
	}

	domain := config.GetConfig().Rollup.SignaturesDomain
	wallet, err := bls.NewWallet(privateKeyDecoded, domain)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

//export NewWalletPrivateKey
func NewWalletPrivateKey() *C.char {
	domain := config.GetConfig().Rollup.SignaturesDomain
	wallet, err := bls.NewRandomWallet(domain)
	if err != nil {
		return nil
	}

	privateKey, _ := wallet.Bytes()
	return C.CString(hex.EncodeToString(privateKey))
}

//export GetWalletPublicKey
func GetWalletPublicKey(privateKey *C.char) *C.char {
	wallet, err := parseWallet(privateKey)
	if err != nil {
		return nil
	}

	publicKey := wallet.PublicKey()
	return C.CString(hex.EncodeToString(publicKey[:]))
}

//export SignTransfer
func SignTransfer(from C.uint, to C.uint, amount C.longlong, fee C.longlong, nonce C.longlong, privateKey *C.char) *C.char {
	wallet, err := parseWallet(privateKey)
	if err != nil {
		return nil
	}

	transfer, _ := api.SignTransfer(wallet, dto.Transfer{
		FromStateID: ref.Uint32(uint32(from)),
		ToStateID:   ref.Uint32(uint32(to)),
		Amount:      models.NewUint256(int64(amount)),
		Fee:         models.NewUint256(int64(fee)),
		Nonce:       models.NewUint256(int64(nonce)),
	})

	return C.CString(hex.EncodeToString(transfer.Signature.Bytes()))
}
