package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestCommander(t *testing.T) {
	cfg := config.GetTestConfig()
	cmd := NewCommander(&cfg)

	err := cmd.Start()
	require.NoError(t, err)

	err = cmd.Stop()
	require.NoError(t, err)
}
