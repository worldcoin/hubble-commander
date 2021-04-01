package bls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignature_Verify(t *testing.T) {
	wallet, err := NewRandomWallet(testDomain)
	require.NoError(t, err)

	data := []byte("0xdeadbeef")
	signature, err := wallet.Sign(data)
	require.NoError(t, err)

	isValid, err := signature.Verify(data, wallet.PublicKey())
	require.NoError(t, err)
	require.True(t, isValid)
}
