package bls

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
