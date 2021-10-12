package bls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	data := []byte{1, 2, 3}

	wallet, err := NewRandomWallet(TestDomain)
	require.NoError(t, err)

	signature, err := wallet.Sign(data)
	require.NoError(t, err)

	privateKey, _ := wallet.Bytes()

	newWallet, err := NewWallet(privateKey, TestDomain)
	require.NoError(t, err)

	isValid, err := signature.Verify(data, newWallet.PublicKey())
	require.NoError(t, err)
	require.True(t, isValid)
}
