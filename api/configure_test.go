package api

import (
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIConfugure(t *testing.T) {
	var enabled bool
	api := API{enableBatchCreation: func(enable bool) {
		enabled = enable
	}}

	err := api.Configure(ConfigureParams{
		CreateBatches: ref.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, true, enabled)

	err = api.Configure(ConfigureParams{
		CreateBatches: ref.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, false, enabled)
}
