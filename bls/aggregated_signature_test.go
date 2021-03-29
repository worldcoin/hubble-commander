package bls

import (
	"testing"

	"github.com/kilic/bn254/bls"
	"github.com/stretchr/testify/require"
)

func TestAggregatedSignature_Verify(t *testing.T) {
	messages := [][]byte{
		[]byte("0x111111"),
		[]byte("0x222222"),
		[]byte("0x333333"),
	}

	publicKeys := make([]*bls.PublicKey, 0, 3)
	signatures := make([]*Signature, 0, 3)

	for _, msg := range messages {
		wallet, err := NewRandomWallet(testDomain)
		require.NoError(t, err)

		sig, err := wallet.Sign(msg)
		require.NoError(t, err)

		publicKeys = append(publicKeys, wallet.PublicKey())
		signatures = append(signatures, sig)
	}

	aggregatedSignature := NewAggregatedSignature(signatures)
	isValid, err := aggregatedSignature.Verify(messages, publicKeys)
	require.NoError(t, err)
	require.True(t, isValid)
}
