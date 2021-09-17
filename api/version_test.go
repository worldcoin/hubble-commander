package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	api := API{cfg: &config.Config{API: &config.APIConfig{Version: "v0123"}}}
	require.Equal(t, "v0123", api.GetVersion())
}
