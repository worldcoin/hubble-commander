package eth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClient_IsActiveProposer(t *testing.T) {
	client, err := NewTestClient()
	require.NoError(t, err)

	isActiveProposer, err := client.IsActiveProposer()
	require.NoError(t, err)

	require.True(t, isActiveProposer)
}
