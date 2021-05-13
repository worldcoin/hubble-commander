package bls

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestDomainFromBytes(t *testing.T) {
	data := crypto.Keccak256([]byte{1, 2, 3})
	domain, err := DomainFromBytes(data)
	require.NoError(t, err)
	require.Equal(t, data, domain.Bytes())
}

func TestDomainFromBytes_InvalidLength(t *testing.T) {
	data := make([]byte, 20)
	domain, err := DomainFromBytes(data)
	require.Equal(t, ErrInvalidDomainLength, err)
	require.Nil(t, domain)
}
