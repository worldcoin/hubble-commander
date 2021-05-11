package bls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWallet(t *testing.T) {
	data := []byte{1, 2, 3}

	wallet, err := NewRandomWallet(testDomain)
	require.NoError(t, err)

	signature, err := wallet.Sign(data)
	require.NoError(t, err)

	privateKey, _ := wallet.Bytes()

	newWallet, err := NewWallet(privateKey, testDomain)
	require.NoError(t, err)

	isValid, err := signature.Verify(data, newWallet.PublicKey())
	require.NoError(t, err)
	require.True(t, isValid)
}

func TestDomainFromBytes(t *testing.T) {
	data := make([]byte, 32)
	domain, err := DomainFromBytes(data)
	require.NoError(t, err)
	require.Equal(t, &Domain{}, domain)
}

func TestDomainFromBytes_InvalidLength(t *testing.T) {
	data := make([]byte, 20)
	domain, err := DomainFromBytes(data)
	require.Equal(t, ErrInvalidDomainLength, err)
	require.Nil(t, domain)
}
