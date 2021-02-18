package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestApi_GetVersion(t *testing.T) {
	api := Api{&config.Config{Version: "v0123"}}
	require.Equal(t, "v0123", api.GetVersion())
}
