package utils

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestParseUnits_OneGWei(t *testing.T) {
	require.Equal(t, ParseUnits("1", params.GWei), big.NewInt(1e9))
}

func TestParseEther_OneEther(t *testing.T) {
	require.Equal(t, ParseEther("1"), big.NewInt(1e18))
}

func TestParseEther_1234Ether(t *testing.T) {
	expected := new(big.Int).Mul(big.NewInt(1234), big.NewInt(1e18))
	require.Equal(t, ParseEther("1234"), expected)
}
