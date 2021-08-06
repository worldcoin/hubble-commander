package commander

import (
	"testing"

	cfg "github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestChainSpec(t *testing.T) {
	config := cfg.GetConfig()
	chainSpec, err := GenerateChainSpec(config)
	require.NoError(t, err)
	require.NotNil(t, *chainSpec)
}
