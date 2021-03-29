package bls

import (
	"testing"

	"github.com/kilic/bn254/bls"
	r "github.com/stretchr/testify/require"
)

func Test_Wallet_SignAndVerify(t *testing.T) {
	bytesToSign := []byte("0x123221")

	wallet, err := NewWallet(DefaultDomain)
	r.NoError(t, err)

	signature, err := wallet.Sign(bytesToSign)
	r.NoError(t, err)

	isValid, err := wallet.VerifySignature(signature, bytesToSign, wallet.signer.Account.Public)
	r.NoError(t, err)
	r.True(t, isValid)
}

func Test_Wallet_SignAndVerifyAggregatedSignature(t *testing.T) {
	totalSigners := 3
	bytesToSign := [3][]byte{
		[]byte("0x111111"),
		[]byte("0x222222"),
		[]byte("0x333333"),
	}

	publicKeys := make([]*bls.PublicKey, 0)
	messages := make([]bls.Message, 0)
	signatures := make([]*bls.Signature, 0)

	for i := 0; i < totalSigners; i++ {
		wallet, err := NewWallet(DefaultDomain)
		r.NoError(t, err)

		signature, err := wallet.Sign(bytesToSign[i])
		r.NoError(t, err)

		publicKeys = append(publicKeys, wallet.signer.Account.Public)
		messages = append(messages, bytesToSign[i])
		signatures = append(signatures, signature)
	}

	aggregatedSignature := NewAggregateSignature(signatures)

	isValid, err := VerifyAggregatedSignature(aggregatedSignature, messages, publicKeys, DefaultDomain)
	r.NoError(t, err)
	r.True(t, isValid)
}
