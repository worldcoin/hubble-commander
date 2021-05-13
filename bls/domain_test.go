package bls

import (
	"encoding/json"
	"fmt"
	"strings"
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

func TestMarshalText(t *testing.T) {
	domain := Domain{1, 2, 3}
	expected := fmt.Sprintf("\"0x010203%s\"", strings.Repeat("0", DomainLength*2-6))
	data, err := json.Marshal(&domain)
	require.NoError(t, err)
	require.Equal(t, expected, string(data))
}
